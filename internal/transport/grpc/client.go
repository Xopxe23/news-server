package grpc_client

import (
	"context"
	"fmt"

	audit "github.com/xopxe23/auditlog/pkg/domain"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Client struct {
	conn        *grpc.ClientConn
	auditClient audit.AuditServiceClient
}

func NewClient(port int) (*Client, error) {
	var conn *grpc.ClientConn
	addr := fmt.Sprintf(":%d", port)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn, auditClient: audit.NewAuditServiceClient(conn)}, nil
}

func (c *Client) CloseConnection() error {
	return c.conn.Close()
}


func (c *Client) SendLogRequest(ctx context.Context, req audit.LogItem) error {
	action, err := audit.ToPbAction(req.Action)
	if err != nil {
		return err
	}

	entity, err := audit.ToPbEntity(req.Entity)
	if err != nil {
		return err
	}

	_, err = c.auditClient.Log(ctx, &audit.LogRequest{
		Action:    action,
		Entity:    entity,
		EntityId:  req.EntityID,
		Timestamp: timestamppb.New(req.Timestamp),
	})

	return err
}