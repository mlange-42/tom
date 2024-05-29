package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

const DefaultUserAgent = "TOM-terminal-open-meteo"

type API string

const (
	OpenMeteo API = "https://api.open-meteo.com/v1"
	Geocoding API = "https://geocoding-api.open-meteo.com/v1"
)

type Client interface {
	Get(ctx context.Context, req Options) ([]byte, error)
}

type openMeteoClient struct {
	URL       string
	UserAgent string
	Client    *http.Client
}

func NewClient(api API) Client {
	return &openMeteoClient{
		URL:       string(api),
		UserAgent: DefaultUserAgent,
		Client:    http.DefaultClient,
	}
}

func (c *openMeteoClient) Get(ctx context.Context, opt Options) ([]byte, error) {
	url := opt.ToURL(c.URL)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("%s - %s", res.Status, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
