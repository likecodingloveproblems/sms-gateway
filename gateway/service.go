package gateway

import (
	"context"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/likecodingloveproblems/sms-gateway/adapter/accounting"
	"github.com/likecodingloveproblems/sms-gateway/entity"
	errmsg "github.com/likecodingloveproblems/sms-gateway/pkg/err_msg"
	"github.com/likecodingloveproblems/sms-gateway/pkg/grpc"
	"github.com/likecodingloveproblems/sms-gateway/types"
	"log"
	"os"
	"strconv"
)

type Service interface {
	ProcessMessage(ctx context.Context, message entity.Message) error
	// This method first debit user's account balance
	// Then add message to the broker to be processed by scheduler
}

type SyncAccountingCheckService struct {
	repo        *Repository
	idGenerator *snowflake.Node
}

func NewService(repository *Repository) Service {
	NodeId := os.Getenv("NODE_ID")
	nodeID, err := strconv.Atoi(NodeId)
	if err != nil {
		log.Fatalf("NODE_ID is necessary: %s", err.Error())
	}
	node, err := snowflake.NewNode(int64(nodeID))
	if err != nil {
		log.Fatalf("Error in snowflake init: %s", err.Error())
	}
	return SyncAccountingCheckService{
		repo:        repository,
		idGenerator: node,
	}
}

func (s SyncAccountingCheckService) calculateMessagePrice(message entity.Message) uint64 {
	// It's assumed that all messages has one price
	return s.repo.GetMessageUnitPrice()
}

func (s SyncAccountingCheckService) ProcessMessage(ctx context.Context, message entity.Message) error {
	// Assign an id to message
	message.ID = types.ID(s.idGenerator.Generate())
	// First calculate its price
	messageCost := s.calculateMessagePrice(message)
	// request from accounting service to reduce users account balance
	err := s.reserveDebit(ctx, message.UserID, messageCost)
	if err != nil {
		// here error can be 500 or `not-enough-budget`
		return err
	}
	// on success add message to be process by scheduler
	return s.processMessage(message)
}

func (s SyncAccountingCheckService) processMessage(message entity.Message) error {
	// Add message to the scheduer stream
	switch message.Type {
	case entity.ExpressMessage:
		err := s.repo.AddMessage(ExpressStreamKey, message)
		if err != nil {
			return err
		}
	case entity.NormalMessage:
		streamKey := fmt.Sprintf(NormalStreamKeyTemplate, message.UserID)
		err := s.repo.AddMessage(streamKey, message)
		if err != nil {
			return err
		}
	default:
		return errmsg.ErrUnknownMessageType
	}
	err := s.repo.AddMessage(ReportingStreamKey, message)
	if err != nil {
		return err
	}
	return nil
}

func (s SyncAccountingCheckService) reserveDebit(ctx context.Context, userId types.ID, cost uint64) error {
	// User the adapter that create a gRPC client to communicate with AccountingService
	conn, err := grpc.NewClient(s.getAccountingGRPCServerConf())
	if err != nil {
		return err
	}
	client := accounting.New(conn)
	if err := client.ReserveDebit(ctx, userId, cost); err != nil {
		return err
	}
	return nil
}

func (s SyncAccountingCheckService) getAccountingGRPCServerConf() grpc.Client {
	panic("not implemented")
}
