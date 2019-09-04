package grpc

import (
	"context"
	"github.com/lucperkins/strato/internal/config"

	"github.com/lucperkins/strato/internal/services/kv"

	"github.com/lucperkins/strato"

	"github.com/lucperkins/strato/proto"

	"google.golang.org/grpc"
)

type GrpcClient struct {
	cacheClient   proto.CacheClient
	counterClient proto.CounterClient
	kvClient      proto.KVClient
	setClient     proto.SetClient
	conn          *grpc.ClientConn
	ctx           context.Context
}

func NewClient(cfg *config.ClientConfig) (*GrpcClient, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	conn, err := connect(cfg.Address)
	if err != nil {
		return nil, err
	}

	cacheClient := proto.NewCacheClient(conn)

	counterClient := proto.NewCounterClient(conn)

	kvClient := proto.NewKVClient(conn)

	setClient := proto.NewSetClient(conn)

	ctx := context.Background()

	return &GrpcClient{
		cacheClient:   cacheClient,
		counterClient: counterClient,
		kvClient:      kvClient,
		setClient:     setClient,
		conn:          conn,
		ctx:           ctx,
	}, nil
}

func connect(addr string) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, grpc.WithInsecure())
}

func (c *GrpcClient) CacheGet(key string) (string, error) {
	req := &proto.CacheGetRequest{
		Key: key,
	}

	val, err := c.cacheClient.CacheGet(c.ctx, req)
	if err != nil {
		return "", err
	}

	return val.Value, nil
}

func (c *GrpcClient) CacheSet(key, value string, ttl int32) error {
	req := &proto.CacheSetRequest{
		Key: key,
		Item: &proto.CacheItem{
			Value: value,
			Ttl:   ttl,
		},
	}

	if _, err := c.cacheClient.CacheSet(c.ctx, req); err != nil {
		return err
	}

	return nil
}

func (c *GrpcClient) IncrementCounter(key string, amount int64) error {
	req := &proto.IncrementCounterRequest{
		Key:    key,
		Amount: amount,
	}

	if _, err := c.counterClient.IncrementCounter(c.ctx, req); err != nil {
		return err
	}

	return nil
}

func (c *GrpcClient) GetCounter(key string) (int64, error) {
	req := &proto.GetCounterRequest{
		Key: key,
	}

	res, err := c.counterClient.GetCounter(c.ctx, req)
	if err != nil {
		return 0, err
	}

	return res.Value, nil
}

func (c *GrpcClient) KVGet(key string) (*kv.Value, error) {
	if key == "" {
		return nil, strato.ErrNoKey
	}

	loc := &proto.Location{
		Key: key,
	}

	res, err := c.kvClient.KVGet(c.ctx, loc)
	if err != nil {
		return nil, err
	}

	val := &kv.Value{
		Content: res.Value.Content,
	}

	return val, nil
}

func (c *GrpcClient) KVPut(key string, value *kv.Value) error {
	if key == "" {
		return strato.ErrNoKey
	}

	if value == nil {
		return strato.ErrNoValue
	}

	loc := &proto.Location{
		Key: key,
	}

	req := &proto.PutRequest{
		Location: loc,
		Value:    value.Proto(),
	}

	if _, err := c.kvClient.KVPut(c.ctx, req); err != nil {
		return err
	}

	return nil
}

func (c *GrpcClient) KVDelete(key string) error {
	if key == "" {
		return strato.ErrNoKey
	}

	loc := &proto.Location{
		Key: key,
	}

	if _, err := c.kvClient.KVDelete(c.ctx, loc); err != nil {
		return err
	}

	return nil
}

func (c *GrpcClient) GetSet(set string) ([]string, error) {
	req := &proto.GetSetRequest{
		Set: set,
	}

	res, err := c.setClient.GetSet(c.ctx, req)
	if err != nil {
		return nil, err
	}

	return res.Items, nil
}

func (c *GrpcClient) AddToSet(set, item string) error {
	req := &proto.ModifySetRequest{
		Set:  set,
		Item: item,
	}

	if _, err := c.setClient.AddToSet(c.ctx, req); err != nil {
		return err
	}

	return nil
}

func (c *GrpcClient) RemoveFromSet(set, item string) error {
	req := &proto.ModifySetRequest{
		Set:  set,
		Item: item,
	}

	if _, err := c.setClient.RemoveFromSet(c.ctx, req); err != nil {
		return err
	}

	return nil
}
