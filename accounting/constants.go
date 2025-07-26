package accounting

import (
	"os"
	"time"
)

const (
	BALANCE_TTL       = time.Minute * 5
	SYNC_CACHE_PERIOD = BALANCE_TTL / 2
	GPRC_SERVER_PORT  = ":50051"
)

var (
	DB_CONNECTION    = os.Getenv("DB")
	REDIS_CONNECTION = os.Getenv("REDIS")
)
