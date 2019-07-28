package strato

import (
	"context"

	"github.com/lucperkins/strato/proto"

	"google.golang.org/grpc"
)

type Client struct {
	cacheClient   proto.CacheClient
	counterClient proto.CounterClient
	kvClient      proto.KVClient
	searchClient  proto.SearchClient
	setClient     proto.SetClient
	conn          *grpc.ClientConn
	ctx           context.Context
}

func NewClient(cfg *ClientConfig) (*Client, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	conn, err := connect(cfg.Address)
	if err != nil {
		return nil, err
	}

	cacheClient := proto.NewCacheClient(conn)

	counterClient := proto.NewCounterClient(conn)

	kvClient := proto.NewKVClient(conn)

	searchClient := proto.NewSearchClient(conn)

	setClient := proto.NewSetClient(conn)

	ctx := context.Background()

	return &Client{
		cacheClient:   cacheClient,
		counterClient: counterClient,
		kvClient:      kvClient,
		searchClient:  searchClient,
		setClient:     setClient,
		conn:          conn,
		ctx:           ctx,
	}, nil
}

func connect(addr string) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, grpc.WithInsecure())
}

func (c *Client) CacheGet(key string) (string, error) {
	req := &proto.CacheGetRequest{
		Key: key,
	}

	val, err := c.cacheClient.CacheGet(c.ctx, req)
	if err != nil {
		return "", err
	}

	return val.Value, nil
}

func (c *Client) CacheSet(key, value string, ttl int32) error {
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

func (c *Client) IncrementCounter(key string, amount int32) error {
	req := &proto.IncrementCounterRequest{
		Key:    key,
		Amount: amount,
	}

	if _, err := c.counterClient.IncrementCounter(c.ctx, req); err != nil {
		return err
	}

	return nil
}

func (c *Client) GetCounter(key string) (int32, error) {
	req := &proto.GetCounterRequest{
		Key: key,
	}

	res, err := c.counterClient.GetCounter(c.ctx, req)
	if err != nil {
		return 0, err
	}

	return res.Value, nil
}

func (c *Client) KVGet(location *Location) (*Value, error) {
	if location == nil {
		return nil, ErrNoLocation
	}

	if err := location.validate(); err != nil {
		return nil, err
	}

	res, err := c.kvClient.Get(c.ctx, location.Proto())
	if err != nil {
		return nil, err
	}

	val := &Value{
		Content: res.Value.Content,
	}

	return val, nil
}

func (c *Client) KVPut(location *Location, value *Value) error {
	if location == nil {
		return ErrNoLocation
	}

	if err := location.validate(); err != nil {
		return err
	}

	if value == nil {
		return ErrNoValue
	}

	req := &proto.PutRequest{
		Location: location.Proto(),
		Value:    value.Proto(),
	}

	if _, err := c.kvClient.KVPut(c.ctx, req); err != nil {
		return err
	}

	return nil
}

func (c *Client) KVDelete(location *Location) error {
	if location == nil {
		return ErrNoLocation
	}

	if err := location.validate(); err != nil {
		return err
	}

	if _, err := c.kvClient.Delete(c.ctx, location.Proto()); err != nil {
		return err
	}

	return nil
}

func (c *Client) Index(doc *Document) error {
	req := &proto.IndexRequest{
		Document: doc.toProto(),
	}

	if _, err := c.searchClient.Index(c.ctx, req); err != nil {
		return err
	}

	return nil
}

func (c *Client) Query(q string) ([]*Document, error) {
	query := &proto.SearchQuery{
		Query: q,
	}

	res, err := c.searchClient.Query(c.ctx, query)
	if err != nil {
		return nil, err
	}

	return docsFromProto(res.Documents), nil
}

func (c *Client) GetSet(set string) ([]string, error) {
	req := &proto.GetSetRequest{
		Set: set,
	}

	res, err := c.setClient.GetSet(c.ctx, req)
	if err != nil {
		return nil, err
	}

	return res.Items, nil
}

func (c *Client) AddToSet(set, item string) error {
	req := &proto.ModifySetRequest{
		Set:  set,
		Item: item,
	}

	if _, err := c.setClient.AddToSet(c.ctx, req); err != nil {
		return err
	}

	return nil
}

func (c *Client) RemoveFromSet(set, item string) error {
	req := &proto.ModifySetRequest{
		Set:  set,
		Item: item,
	}

	if _, err := c.setClient.RemoveFromSet(c.ctx, req); err != nil {
		return err
	}

	return nil
}
