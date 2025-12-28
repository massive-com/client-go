package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/massive-com/client-go/v2/rest/encoder"
	"github.com/massive-com/client-go/v2/rest/models"
)

const clientVersion = "v1.16.0"

const (
	APIURL            = "https://api.massive.com"
	DefaultRetryCount = 3
)

// Client defines an HTTP client for the Massive REST API.
type Client struct {
	HTTP    *resty.Client
	encoder *encoder.Encoder
}

// New returns a new client with the specified API key and default settings.
func New(apiKey string) Client {
	return newClient(apiKey, nil)
}

// NewWithClient returns a new client with the specified API key and a custom HTTP client.
func NewWithClient(apiKey string, hc *http.Client) Client {
	return newClient(apiKey, hc)
}

func newClient(apiKey string, hc *http.Client) Client {
	var c *resty.Client
	if hc == nil {
		c = resty.New()
	} else {
		c = resty.NewWithClient(hc)
	}

	c.SetBaseURL(APIURL)
	c.SetAuthToken(apiKey)
	c.SetRetryCount(DefaultRetryCount)
	c.SetTimeout(10 * time.Second)
	c.SetHeader("User-Agent", fmt.Sprintf("Massive.com GoClient/%v", clientVersion))
	c.SetHeader("Accept-Encoding", "gzip")

	return Client{
		HTTP:    c,
		encoder: encoder.New(),
	}
}

// Call makes an API call based on the request params and options. The response is automatically unmarshaled.
func (c *Client) Call(ctx context.Context, method, path string, params, response any, opts ...models.RequestOption) error {
	uri, err := c.encoder.EncodeParams(path, params)
	if err != nil {
		return err
	}
	return c.CallURL(ctx, method, uri, response, opts...)
}

// CallURL makes an API call based on a request URI and options. The response is automatically unmarshaled.
func (c *Client) CallURL(ctx context.Context, method, uri string, response any, opts ...models.RequestOption) error {
	options := mergeOptions(opts...)

	req := c.HTTP.R().SetContext(ctx)
	if options.APIKey != nil {
		req.SetAuthToken(*options.APIKey)
	}
	req.SetQueryParamsFromValues(options.QueryParams)
	req.SetHeaderMultiValues(options.Headers)
	req.SetResult(response).SetError(&models.ErrorResponse{})

	if options.Limiter != nil {
		options.Limiter.Take()
	}

	res, err := req.Execute(method, uri)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	} else if res.IsError() {
		errRes := res.Error().(*models.ErrorResponse)
		errRes.StatusCode = res.StatusCode()
		if errRes.RequestID == "" {
			errRes.RequestID = res.Header().Get("X-Request-ID")
		}
		return errRes
	}

	if options.Trace {
		fmt.Printf("Request URL: %s\n", uri)
		sanitizedHeaders := req.Header
		for k := range sanitizedHeaders {
			if k == "Authorization" {
				sanitizedHeaders[k] = []string{"REDACTED"}
			}
		}
		fmt.Printf("Request Headers: %s\n", sanitizedHeaders)
		fmt.Printf("Response Headers: %+v\n", res.Header())
	}

	return nil
}

func mergeOptions(opts ...models.RequestOption) *models.RequestOptions {
	options := &models.RequestOptions{}
	for _, o := range opts {
		o(options)
	}

	return options
}
