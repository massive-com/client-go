// Package massive defines a REST client for the Massive API.
package massive

import (
	"net/http"

	"github.com/massive-com/client-go/v2/rest/client"
)

// Client defines a client to the Massive REST API.
type Client struct {
	client.Client
	AggsClient
	QuotesClient
	ReferenceClient
	TradesClient
	SnapshotClient
	IndicatorsClient
	SummariesClient
	FuturesClient
	VX VXClient
}

// New creates a client for the Massive REST API.
func New(apiKey string) *Client {
	return newClient(apiKey, nil)
}

// NewWithClient creates a client for the Massive REST API using a custom HTTP client.
func NewWithClient(apiKey string, hc *http.Client) *Client {
	return newClient(apiKey, hc)
}

func newClient(apiKey string, hc *http.Client) *Client {
	var c client.Client
	if hc == nil {
		c = client.New(apiKey)
	} else {
		c = client.NewWithClient(apiKey, hc)
	}

	return &Client{
		Client:           c,
		IndicatorsClient: IndicatorsClient{Client: c},
		SummariesClient:  SummariesClient{Client: c},
		AggsClient:       AggsClient{Client: c},
		QuotesClient:     QuotesClient{Client: c},
		ReferenceClient:  ReferenceClient{Client: c},
		TradesClient:     TradesClient{Client: c},
		SnapshotClient:   SnapshotClient{Client: c},
		FuturesClient:    FuturesClient{Client: c},
		VX:               VXClient{Client: c},
	}
}
