package strato

import (
	"context"
	"strato/proto"

	"google.golang.org/grpc"
)

type Client struct {
	kvClient proto.KVClient
	conn     *grpc.ClientConn
	ctx      context.Context
}

func NewClient(cfg *ClientConfig) (*Client, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	conn, err := connect(cfg.Address)
	if err != nil {
		return nil, err
	}

	kvClient := proto.NewKVClient(conn)

	ctx := context.Background()

	return &Client{
		kvClient: kvClient,
		conn:     conn,
		ctx:      ctx,
	}, nil
}

func connect(addr string) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, grpc.WithInsecure())
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
