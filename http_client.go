package strato

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/lucperkins/strato/internal/services/cache"
	"net/http"
)

type HttpClient struct {
	rootUrl string
	cl      *resty.Client
}

var _ cache.Cache = (*HttpClient)(nil)

func NewHttpClient(cfg *ClientConfig) *HttpClient {
	cl := resty.New()

	return &HttpClient{
		rootUrl: cfg.Address,
		cl:      cl,
	}
}

func (c *HttpClient) cacheKeyUrl(key string) string {
	return fmt.Sprintf("%s/cache/%s", c.rootUrl, key)
}

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
			"ttl":   string(ttl),
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
