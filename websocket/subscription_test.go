package massivews

import (
	"encoding/json"
	"testing"

	"github.com/massive-com/client-go/v3/websocket/models"
	"github.com/stretchr/testify/assert"
)

func TestSupportsTopic(t *testing.T) {
	assert.Equal(t, true, Stocks.supports(StocksMinAggs))
	assert.Equal(t, false, Stocks.supports(stocksMax))
	assert.Equal(t, true, Options.supports(OptionsSecAggs))
	assert.Equal(t, false, Options.supports(StocksMinAggs))
	assert.Equal(t, true, Forex.supports(ForexQuotes))
	assert.Equal(t, false, Forex.supports(OptionsQuotes))
	assert.Equal(t, true, Crypto.supports(CryptoL2Book))
	assert.Equal(t, false, Crypto.supports(cryptoMin))
	assert.Equal(t, true, Market("testMarket").supports(StocksImbalances))
}

func TestGet(t *testing.T) {
	c, err := New(Config{
		APIKey: "test",
		Feed:   RealTime,
		Market: Stocks,
	})
	assert.Nil(t, err)

	// subscribe to a topic
	minAggs := models.ControlMessage{Action: models.Subscribe, Params: "AM.AAPL"}
	aBytes, _ := json.Marshal(minAggs)
	err = c.Subscribe(StocksMinAggs, "AAPL")
	assert.Nil(t, err)
	msgs := c.subs.get()
	assert.Equal(t, []json.RawMessage{aBytes}, msgs)

	// unsubscribe from *
	err = c.Unsubscribe(StocksMinAggs, "NFLX", "*")
	assert.Nil(t, err)
	msgs = c.subs.get()
	assert.Equal(t, []json.RawMessage(nil), msgs)

	// subscribe to another topic
	err = c.Subscribe(StocksTrades, "SNAP")
	assert.Nil(t, err)
	trades := models.ControlMessage{Action: models.Subscribe, Params: "T.SNAP"}
	tBytes, _ := json.Marshal(trades)
	msgs = c.subs.get()
	assert.Equal(t, []json.RawMessage{tBytes}, msgs)

	// unsubscribe from *
	err = c.Unsubscribe(StocksTrades)
	assert.Nil(t, err)
	msgs = c.subs.get()
	assert.Equal(t, []json.RawMessage(nil), msgs)
}

func TestGetSub(t *testing.T) {
	submsg, err := getSub(models.Subscribe, StocksMinAggs, "AAPL", "GME", "HOOD")
	assert.Nil(t, err)
	sub, err := json.Marshal(models.ControlMessage{Action: models.Subscribe, Params: "AM.AAPL,AM.GME,AM.HOOD"})
	assert.Nil(t, err)
	assert.Equal(t, json.RawMessage(sub), submsg)

	unsubmsg, err := getSub(models.Unsubscribe, StocksSecAggs)
	assert.Nil(t, err)
	unsub, err := json.Marshal(models.ControlMessage{Action: models.Unsubscribe, Params: "A.*"})
	assert.Nil(t, err)
	assert.Equal(t, json.RawMessage(unsub), unsubmsg)
}

func TestSubscriptions(t *testing.T) {
	c, err := New(Config{
		APIKey: "test",
		Feed:   RealTime,
		Market: Stocks,
	})
	assert.Nil(t, err)

	err = c.Subscribe(StocksMinAggs, "AAPL", "TSLA")
	assert.Nil(t, err)
	_, aapl := c.subs[StocksMinAggs]["AAPL"]
	assert.True(t, aapl)
	_, tsla := c.subs[StocksMinAggs]["TSLA"]
	assert.True(t, tsla)

	err = c.Unsubscribe(StocksMinAggs, "AAPL", "NFLX")
	assert.Nil(t, err)
	_, aapl = c.subs[StocksMinAggs]["AAPL"]
	assert.False(t, aapl)

	err = c.Subscribe(StocksMinAggs)
	assert.Nil(t, err)
	_, all := c.subs[StocksMinAggs]["*"]
	assert.True(t, all)
	_, tsla = c.subs[StocksMinAggs]["TSLA"]
	assert.False(t, tsla)

	err = c.Unsubscribe(StocksMinAggs, "*")
	assert.Nil(t, err)
	_, all = c.subs[StocksMinAggs]["*"]
	assert.False(t, all)
	_, trade := c.subs[StocksTrades]
	assert.False(t, trade)

	err = c.Unsubscribe(StocksTrades, "RDFN")
	assert.Nil(t, err)
	_, trade = c.subs[StocksTrades]
	assert.False(t, trade)

	err = c.Subscribe(StocksTrades, "FB")
	assert.Nil(t, err)
	_, fb := c.subs[StocksTrades]["FB"]
	assert.True(t, fb)

	err = c.Unsubscribe(StocksTrades)
	assert.Nil(t, err)
	_, fb = c.subs[StocksTrades]["FB"]
	assert.False(t, fb)
}
