package accounting

import (
	"context"
	"database/sql"
	reservedebit "github.com/likecodingloveproblems/sms-gateway/protobuf/accounting/golang/ReserveDebit"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	reservedebit.UnimplementedAccountingServer
	storage Storage
}

func (s *server) ReserveDebit(ctx context.Context, req *reservedebit.ReserveDebitRequest) (*reservedebit.ReserveDebitResponse, error) {
}

func (s *server) CancelDebit(ctx context.Context, req *reservedebit.CancelDebitRequest) (*reservedebit.CancelDebitResponse, error) {

}

func Serve() {
	var err error
	ctx := context.Background()
	db, err := sql.Open("postgres", DB_CONNECTION)
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: REDIS_CONNECTION,
		DB:   0,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis error: %v", err)
	}

	storage := NewCachedDB(ctx, db, rdb)
	gRPCServer := &server{
		storage: storage,
	}

	lis, err := net.Listen("tcp", GPRC_SERVER_PORT)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	reservedebit.RegisterAccountingServer(s, gRPCServer)
	log.Println("Accounting gRPC server running at :50051")
	err = s.Serve(lis)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}
