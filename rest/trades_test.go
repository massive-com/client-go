package massive_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	massive "github.com/massive-com/client-go/v2/rest"
	"github.com/massive-com/client-go/v2/rest/models"
	"github.com/stretchr/testify/assert"
)

func TestListTrades(t *testing.T) {
	c := massive.New("API_KEY")

	httpmock.ActivateNonDefault(c.HTTP.GetClient())
	defer httpmock.DeactivateAndReset()

	trade1 := `{
	"conditions": [
		12,
		41
	],
	"exchange": 11,
	"id": "1",
	"participant_timestamp": 1517562000015577000,
	"price": 171.55,
	"sequence_number": 1063,
	"sip_timestamp": 1517562000016036600,
	"size": 100,
	"decimal_size": "100.0",
	"tape": 3
}`

	trade2 := `{
	"conditions": [
		12,
		41
	],
	"exchange": 11,
	"id": "2",
	"participant_timestamp": 1517562000015577600,
	"price": 171.55,
	"sequence_number": 1064,
	"sip_timestamp": 1517562000016038100,
	"size": 100,
	"decimal_size": "100.233456",
	"tape": 3
}`

	expectedResponse := `{
	"status": "OK",
	"request_id": "a47d1beb8c11b6ae897ab76cdbbf35a3",
	"next_url": "https://api.massive.com/v3/trades/AAPL?cursor=YWN0aXZlPXRydWUmZGF0ZT0yMDIxLTA0LTI1JmxpbWl0PTEmb3JkZXI9YXNjJnBhZ2VfbWFya2VyPUElN0M5YWRjMjY0ZTgyM2E1ZjBiOGUyNDc5YmZiOGE1YmYwNDVkYzU0YjgwMDcyMWE2YmI1ZjBjMjQwMjU4MjFmNGZiJnNvcnQ9dGlja2Vy",
	"results": [
` + indent(true, trade1, "\t\t") + `,
` + indent(true, trade2, "\t\t") + `
	]
}`

	registerResponder("https://api.massive.com/v3/trades/AAPL?limit=2&order=asc&sort=timestamp&timestamp.gte=1626948000000000000", expectedResponse)
	registerResponder("https://api.massive.com/v3/trades/AAPL?cursor=YWN0aXZlPXRydWUmZGF0ZT0yMDIxLTA0LTI1JmxpbWl0PTEmb3JkZXI9YXNjJnBhZ2VfbWFya2VyPUElN0M5YWRjMjY0ZTgyM2E1ZjBiOGUyNDc5YmZiOGE1YmYwNDVkYzU0YjgwMDcyMWE2YmI1ZjBjMjQwMjU4MjFmNGZiJnNvcnQ9dGlja2Vy&sort=timestamp", "{}")
	iter := c.ListTrades(context.Background(), models.ListTradesParams{Ticker: "AAPL"}.
		WithTimestamp(models.GTE, models.Nanos(time.Date(2021, 7, 22, 10, 0, 0, 0, time.UTC))).
		WithOrder(models.Asc).WithLimit(2), models.QueryParam("sort", string(models.Timestamp)))

	// iter creation
	assert.Nil(t, iter.Err())
	assert.NotNil(t, iter.Item())

	// first item
	assert.True(t, iter.Next())
	assert.Nil(t, iter.Err())
	var expect1 models.Trade
	err := json.Unmarshal([]byte(trade1), &expect1)
	assert.Nil(t, err)
	assert.Equal(t, expect1, iter.Item())

	// second item
	assert.True(t, iter.Next())
	assert.Nil(t, iter.Err())
	var expect2 models.Trade
	err = json.Unmarshal([]byte(trade2), &expect2)
	assert.Nil(t, err)
	assert.Equal(t, expect2, iter.Item())

	// end of list
	assert.False(t, iter.Next())
	assert.Nil(t, iter.Err())
}

func TestGetLastTrade(t *testing.T) {
	c := massive.New("API_KEY")

	httpmock.ActivateNonDefault(c.HTTP.GetClient())
	defer httpmock.DeactivateAndReset()

	expectedResponse := `{
	"status": "OK",
	"request_id": "f05562305bd26ced64b98ed68b3c5d96",
	"results": {
		"T": "AAPL",
		"f": 1617901342969796400,
		"q": 3135876,
		"t": 1617901342969834000,
		"y": 1617901342968000000,
		"c": [
			37
		],
		"i": "118749",
		"p": 129.8473,
		"r": 202,
		"s": 25,
		"x": 4,
		"z": 3
	}
}`

	registerResponder("https://api.massive.com/v2/last/trade/AAPL", expectedResponse)
	res, err := c.GetLastTrade(context.Background(), &models.GetLastTradeParams{
		Ticker: "AAPL",
	})
	assert.Nil(t, err)

	var expect models.GetLastTradeResponse
	err = json.Unmarshal([]byte(expectedResponse), &expect)
	assert.Nil(t, err)
	assert.Equal(t, &expect, res)
}

func TestGetLastCryptoTrade(t *testing.T) {
	c := massive.New("API_KEY")

	httpmock.ActivateNonDefault(c.HTTP.GetClient())
	defer httpmock.DeactivateAndReset()

	expectedResponse := `{
	"status": "success",
	"request_id": "d2d779df015fe2b7fbb8e58366610ef7",
	"symbol": "BTC-USD",
	"last": {
		"price": 16835.42,
		"size": 0.006909,
		"exchange": 4,
		"conditions": [
			1
		],
		"timestamp": 1605560885027
	}
}`

	registerResponder("https://api.massive.com/v1/last/crypto/BTC/USD", expectedResponse)
	res, err := c.GetLastCryptoTrade(context.Background(), &models.GetLastCryptoTradeParams{
		From: "BTC",
		To:   "USD",
	})
	assert.Nil(t, err)

	var expect models.GetLastCryptoTradeResponse
	err = json.Unmarshal([]byte(expectedResponse), &expect)
	assert.Nil(t, err)
	assert.Equal(t, &expect, res)
}
