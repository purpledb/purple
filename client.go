package strato

import (
	"context"
	"google.golang.org/grpc"
	"strato/proto"
)

type Client struct {
	kvClient proto.KVClient
	ctx      context.Context
}

var _ KV = (*Client)(nil)

func NewClient(cfg *ClientConfig) (*Client, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(cfg.Address)
	if err != nil {
		return nil, err
	}

	kvClient := proto.NewKVClient(conn)

	ctx := context.Background()

	return &Client{
		kvClient: kvClient,
		ctx:      ctx,
	}, nil
}

func (c *Client) Get(location *Location) (*Value, error) {
	res, err := c.kvClient.Get(c.ctx, location.Proto())
	if err != nil {
		return nil, err
	}

	val := &Value{
		Content: res.Value.Content,
	}

	return val, nil
}

func (c *Client) Put(location *Location, value *Value) error {
	if location == nil {
		return ErrNoLocation
	}

	if value == nil {
		return ErrNoValue
	}

	req := &proto.PutRequest{
		Location: location.Proto(),
		Value:    value.Proto(),
	}

	if _, err := c.kvClient.Put(c.ctx, req); err != nil {
		return err
	}

	return nil
}

func (c *Client) Delete(location *Location) error {
	if _, err := c.kvClient.Delete(c.ctx, location.Proto()); err != nil {
		return err
	}

	return nil
}
