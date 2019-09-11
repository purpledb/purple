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
	_ kv.KV           = (*HttpClient)(nil)
	//_ set.Set         = (*HttpClient)(nil)
)

func NewHttpClient(cfg *ClientConfig) (*HttpClient, error) {
	cl := resty.New()

	client := &HttpClient{
		rootUrl: cfg.Address,
		cl:      cl,
	}

	if err := client.ping(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *HttpClient) ping() error {
	url := c.pingUrl()

	res, err := c.cl.R().
		Get(url)

	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusOK {
		return ErrHttpUnavailable
	}

	return nil
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
		return fmt.Errorf("expected status code 204, got %d", res.StatusCode())
	}

	return nil
}

func (c *HttpClient) KVDelete(key string) error {
	url := c.kvKeyUrl(key)

	res, err := c.cl.R().
		Delete(url)

	if err != nil {
		return err
	}

	if res.StatusCode() != http.StatusNoContent {
		return fmt.Errorf("expected status code 204, got %d", res.StatusCode())
	}

	return nil
}

type setResponse struct {
	Items []string `json:"items"`
	Set   string   `json:"set"`
}

// Set
func (c *HttpClient) SetGet(key string) ([]string, error) {
	var js setResponse

	url := c.setKeyUrl(key)

	res, err := c.cl.R().
		Get(url)

	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("expected status code 200, got %d", res.StatusCode())
	}

	if err := json.Unmarshal(res.Body(), &js); err != nil {
		return nil, err
	}

	return js.Items, nil
}

func (c *HttpClient) SetAdd(key, item string) ([]string, error) {
	url := c.setKeyUrl(key)

	var js setResponse

	res, err := c.cl.R().
		SetQueryParam("item", item).
		Put(url)

	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("expected status code 200, got %d", res.StatusCode())
	}

	if err := json.Unmarshal(res.Body(), &js); err != nil {
		return nil, err
	}

	return js.Items, nil
}

func (c *HttpClient) SetRemove(key, item string) ([]string, error) {
	url := c.setKeyUrl(key)

	var js setResponse

	res, err := c.cl.R().
		SetQueryParam("item", item).
		Delete(url)

	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("expected status code 200, got %d", res.StatusCode())
	}

	if err := json.Unmarshal(res.Body(), &js); err != nil {
		return nil, err
	}

	return js.Items, nil
}

// Helpers
func (c *HttpClient) pingUrl() string {
	return fmt.Sprintf("%s/ping", c.rootUrl)
}

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

func (c *HttpClient) setKeyUrl(key string) string {
	return keyUrl(c.rootUrl, "sets", key)
}

func int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

func int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
