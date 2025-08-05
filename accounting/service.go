package accounting

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/likecodingloveproblems/sms-gateway/pkg/cron_job"
	errmsg "github.com/likecodingloveproblems/sms-gateway/pkg/err_msg"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"strings"
)

type Storage interface {
	Withdraw(ctx context.Context, userId uint64, amount uint64) (isSuccessful bool, err error)
	Deposit(ctx context.Context, userId uint64, amount uint64) (isSuccessful bool, err error)
}

type CachedPostgresStorage struct {
	rdb                    *redis.Client
	db                     *sql.DB
	syncDBWithCacheCronJob *cron_job.CronJob
}

func NewCachedDB(ctx context.Context, db *sql.DB, rdb *redis.Client) *CachedPostgresStorage {
	storage := &CachedPostgresStorage{
		db:  db,
		rdb: rdb,
	}
	syncCacheJob := cron_job.NewCronJob(ctx, SYNC_CACHE_PERIOD, storage.syncDBWithCache)
	storage.syncDBWithCacheCronJob = syncCacheJob
	return storage
}

func (c *CachedPostgresStorage) SyncDBWithCache(ctx context.Context) {
	go c.syncDBWithCacheCronJob.Run()
}

func (c *CachedPostgresStorage) syncDBWithCache(ctx context.Context) {
	cursor := uint64(0)
	var keys []string
	var err error

	// It's possible that when we fetch balances from redis and go update database, until updating database
	// cache be invalidated and fetched from database, then we have data in consistency

	// SCAN keys matching pattern balance:*
	// Almost in the requirements number of businesses is around 100k
	// So it's not too big to require huge amount of RAM
	for {
		var scannedKeys []string
		scannedKeys, cursor, err = c.rdb.Scan(ctx, cursor, "balance:*", 100).Result()
		if err != nil {
			log.Printf("Redis scan error: %v", err)
			break
		}
		keys = append(keys, scannedKeys...)
		if cursor == 0 {
			break
		}
	}

	if len(keys) == 0 {
		return
	}

	// Get values for keys
	vals, err := c.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		log.Printf("Redis MGET error: %v", err)
		return
	}

	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	for i, key := range keys {
		if vals[i] == nil {
			continue
		}
		userID := strings.TrimPrefix(key, "balance:")
		balanceStr := fmt.Sprintf("%v", vals[i])
		balance, err := strconv.ParseInt(balanceStr, 10, 64)
		if err != nil {
			log.Printf("Invalid balance format for %s: %v", key, err)
			continue
		}
		updates = append(updates, fmt.Sprintf("($%d, $%d)", argIndex, argIndex+1))
		args = append(args, userID, balance)
		argIndex += 2
	}
	tx, err := c.db.Begin()
	if err != nil {
		log.Printf("DB begin error: %v", err)
	}

	if len(updates) > 0 {
		query := `
		UPDATE accounts AS a SET balance = b.balance
		FROM (VALUES ` + strings.Join(updates, ", ") + `
		) AS b(user_id, balance)
		WHERE a.user_id = b.user_id;
	`
		if _, err := tx.Exec(query, args...); err != nil {
			log.Printf("Batch update error: %v", err)
			_ = tx.Rollback()
			return
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("DB commit error: %v", err)
	}
}

func (c *CachedPostgresStorage) Deposit(ctx context.Context, userId uint64, amount uint64) (isSuccessful bool, err error) {
	if err := c.cacheBalanceFromDBIfNecessary(ctx, userId); err != nil {
		return false, err
	}
	if _, err = c.rdb.IncrBy(ctx, c.getUserBalanceKey(userId), int64(amount)).Result(); err != nil {
		return false, err
	}
	return true, nil
}

func (c *CachedPostgresStorage) Withdraw(ctx context.Context, userId uint64, amount uint64) (isSuccessful bool, err error) {
	if err := c.cacheBalanceFromDBIfNecessary(ctx, userId); err != nil {
		return false, err
	}
	isSuccessful, err = c.withdraw(ctx, userId, amount)
	return isSuccessful, err
}

func (c *CachedPostgresStorage) withdraw(ctx context.Context, userId uint64, amount uint64) (bool, error) {
	// Here it's important to handle race condition
	// so if we do it in the redis as an atomic operation we don't have race condition
	luaDecrementScript := redis.NewScript(`
		local key = KEYS[1]
		local amt = tonumber(ARGV[1])
		local bal = tonumber(redis.call("GET", key) or "-1")
		if bal == -1 then return {err="not_found"} end
		if bal < amt then return 0 end
		redis.call("DECRBY", key, amt)
		return 1
	`)
	res, err := luaDecrementScript.Run(ctx, c.rdb, []string{c.getUserBalanceKey(userId)}, amount).Result()
	if err != nil {
		if err.Error() == "not_found" {
			// As we check this condition, it must not be reached any time
			return false, errmsg.ErrUserNotFound
		}
		return false, err
	}
	if res.(int64) == 0 {
		return false, errors.New("insufficient_balance")
	}
	return true, nil
}

func (c *CachedPostgresStorage) balanceExistsInCache(ctx context.Context, userId uint64) (bool, error) {
	exists, err := c.rdb.Exists(ctx, c.getUserBalanceKey(userId)).Result()
	return exists == 1, err
}

func (c *CachedPostgresStorage) cacheBalanceFromDB(ctx context.Context, userId uint64) error {
	// Cache update can lead to stampedes or thundering herd problem.
	// Fetch value from postgres table and cache balance of user in redis
	balance, err := c.fetchBalanceFromDB(ctx, userId)
	if errors.Is(err, sql.ErrNoRows) {
		return errmsg.ErrUserNotFound
	}
	if err != nil {
		return err
	}
	// update cache into the redis
	err = c.storeBalanceInCache(ctx, userId, balance)
	if err != nil {
		return err
	}
	return nil
}

func (c *CachedPostgresStorage) getUserBalanceKey(userId uint64) string {
	return fmt.Sprintf("balance:%d", userId)
}

func (c *CachedPostgresStorage) fetchBalanceFromDB(ctx context.Context, userId uint64) (uint64, error) {
	var balance uint64
	if err := c.db.QueryRow("SELECT balance FROM accounts WHERE user_id = $1", userId).Scan(&balance); err != nil {
		return 0, err
	}
	return balance, nil
}

func (c *CachedPostgresStorage) storeBalanceInCache(ctx context.Context, userId uint64, balance uint64) error {
	return c.rdb.Set(ctx, c.getUserBalanceKey(userId), balance, BALANCE_TTL).Err()
}

func (c *CachedPostgresStorage) cacheBalanceFromDBIfNecessary(ctx context.Context, userId uint64) error {
	exists, err := c.balanceExistsInCache(ctx, userId)
	if err != nil {
		return err
	}
	if !exists {
		err := c.cacheBalanceFromDB(ctx, userId)
		if err != nil {
			return errmsg.ErrUnexpectedError
		}
	}
	return nil
}

var _ Storage = &CachedPostgresStorage{}

func ReserveDebit() {
	// Atomic
	// Decrease users balance
	// instance status pending
	// Failure Cache
}

func CommitDebit() {
	// Status commit
}

func Rollback() {
	// Increase
}

func Expiration() {
	// Expire
	// Increase
}

// ACID
// RDBMS
// Write -> 5000 write per second
// Sync

// Load balancer
// cluster postgres

// Scalable
// 100,000
// 1,000,000
