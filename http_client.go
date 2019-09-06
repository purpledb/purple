package strato

import (
	"encoding/json"
	"fmt"
	"github.com/lucperkins/strato/internal/services/counter"
	"github.com/lucperkins/strato/internal/services/kv"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/lucperkins/strato/internal/services/cache"
)

type HttpClient struct {
	rootUrl string
	cl      *resty.Client
}

var (
	_ cache.Cache     = (*HttpClient)(nil)
	_ counter.Counter = (*HttpClient)(nil)
	//_ kv.KV           = (*HttpClient)(nil)
)

func NewHttpClient(cfg *ClientConfig) *HttpClient {
	cl := resty.New()

	return &HttpClient{
		rootUrl: cfg.Address,
		cl:      cl,
	}
}

// Cache operations
func (c *HttpClient) CacheGet(key string) (string, error) {
	type value struct {
		Value string `json:"value"`
	}

	var val value

	url := c.cacheKeyUrl(key)

	res, err := c.cl.R().Get(url)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(res.Body(), &val); err != nil {
		return "", err
	}

	return val.Value, nil
}

func (c *HttpClient) CacheSet(key, value string, ttl int32) error {
	url := c.cacheKeyUrl(key)

	res, err := c.cl.R().
		SetQueryParams(map[string]string{
			"value": value,
			"ttl":   int32ToString(ttl),
		}).
		Put(url)

	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("expected status code 204, got %d", res.StatusCode())
	}

	return nil
}

// Counter operations
func (c *HttpClient) CounterGet(key string) (int64, error) {
	type value struct {
		Value int64 `json:"value"`
	}

	var val value

	url := c.counterKeyUrl(key)

	res, err := c.cl.R().
		Get(url)

	if err != nil {
		return 0, err
	}

	if res.StatusCode() != http.StatusOK {
		return 0, fmt.Errorf("expected status code 200, got %d", res.StatusCode())
	}

	if err := json.Unmarshal(res.Body(), &val); err != nil {
		return 0, err
	}

	return val.Value, nil
}

func (c *HttpClient) CounterIncrement(key string, increment int64) error {
	url := c.counterKeyUrl(key)

	res, err := c.cl.R().
		SetQueryParam("increment", int64ToString(increment)).
		Put(url)

	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("expected status code 204, got %d", res.StatusCode())
	}

	return nil
}

// KV
func (c *HttpClient) KVGet(key string) (*kv.Value, error) {
	type value struct {
		Value []byte `json:"value"`
	}

	var val value

	url := c.kvKeyUrl(key)

	res, err := c.cl.R().
		Get(url)

	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("expected status code 200, got %d", res.StatusCode())
	}

	if err := json.Unmarshal(res.Body(), &val); err != nil {
		return nil, err
	}

	return &kv.Value{
		Content: val.Value,
	}, nil
}

func (c *HttpClient) KVPut(key string, value *kv.Value) error {
	js := map[string][]byte{
		"content": value.Content,
	}

	url := c.kvKeyUrl(key)

	res, err := c.cl.R().
		SetBody(js).
		Put(url)

	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusNoContent {
		fmt.Println(res.String())
		return fmt.Errorf("expected status code 204, got %d", res.StatusCode())
	}

	return nil
}

// Helpers
func keyUrl(root, service, key string) string {
	return fmt.Sprintf("%s/%s/%s", root, service, key)
}

func (c *HttpClient) cacheKeyUrl(key string) string {
	return keyUrl(c.rootUrl, "cache", key)
}

func (c *HttpClient) counterKeyUrl(key string) string {
	return keyUrl(c.rootUrl, "counters", key)
}

func (c *HttpClient) kvKeyUrl(key string) string {
	return keyUrl(c.rootUrl, "kv", key)
}

func int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

func int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
