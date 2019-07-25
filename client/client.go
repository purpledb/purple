package client

import (
	"context"
	"google.golang.org/grpc"
	"strato/kv"
	"strato/proto"
)

type Client struct {
	kvClient proto.KVClient
	ctx      context.Context
}

var _ kv.KV = (*Client)(nil)

func New(cfg *Config) (*Client, error) {
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

func (c *Client) Get(location *kv.Location) (*kv.Value, error) {
	res, err := c.kvClient.Get(c.ctx, location.Proto())
	if err != nil {
		return nil, err
	}

	val := &kv.Value{
		Content: res.Value.Content,
	}

	return val, nil
}

func (c *Client) Put(location *kv.Location, value *kv.Value) error {
	req := &proto.PutRequest{
		Location: location.Proto(),
		Value: value.Proto(),
	}

	if _, err := c.kvClient.Put(c.ctx, req); err != nil {
		return err
	}

	return nil
}

func (c *Client) Delete(location *kv.Location) error {
	if _, err := c.kvClient.Delete(c.ctx, location.Proto()); err != nil {
		return err
	}

	return nil
}
