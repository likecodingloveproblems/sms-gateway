package accounting

import (
	"context"
	"errors"
	reservedebit "github.com/likecodingloveproblems/sms-gateway/protobuf/accounting/golang/ReserveDebit"
	"github.com/likecodingloveproblems/sms-gateway/types"
	"google.golang.org/grpc"
)

type Client struct {
	conn *grpc.ClientConn
}

func New(conn *grpc.ClientConn) *Client {
	return &Client{
		conn: conn,
	}
}

func (c Client) ReserveDebit(ctx context.Context, userID types.ID, amount uint64) error {
	client := reservedebit.NewAccountingClient(c.conn)
	req := &reservedebit.ReserveDebitRequest{
		UserId: uint64(userID),
		Amount: amount,
	}
	resp, err := client.ReserveDebit(ctx, req)
	if err != nil {
		// Now it must be retried with idempotency
		// After max retry return ErrInternalServer
	}
	if resp.IsSuccessful {
		return nil
	}
	return errors.New(resp.Reason)
}

func (c Client) CancelDebit(ctx context.Context, userID types.ID, amount uint64) (bool, string) {
	client := reservedebit.NewAccountingClient(c.conn)
	req := &reservedebit.CancelDebitRequest{
		UserId: uint64(userID),
		Amount: amount,
	}
	resp, err := client.CancelDebit(ctx, req)
	if err != nil {
		// Now it must be retried with idempotency
		// After max retry return ErrInternalServer
	}
	if resp.IsSuccessful {
		return true, ""
	}
	return false, resp.Reason
}
