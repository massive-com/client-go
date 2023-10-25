package polygonws

import (
	"errors"
)

// Config is a set of WebSocket client options.
type Config struct {
	// APIKey is the API key used to authenticate against the server.
	APIKey string

	// Feed is the data feed (e.g. Delayed, RealTime) which represents the server host.
	Feed Feed

	// Market is the type of market (e.g. Stocks, Crypto) used to connect to the server.
	Market Market

	// MaxRetries is the maximum number of retry attempts that will occur. If the maximum
	// is reached, the client will close the connection. Omitting this will cause the
	// client to reconnect indefinitely until the user closes it.
	MaxRetries *uint64

	// RawData is a flag indicating whether data should be returned as a raw JSON or raw bytes. If BypassRawDataRouting is unset
	// then the data will be returned as raw JSON, otherwise it will be raw bytes.
	RawData bool

	// BypassRawDataRouting is a flag that interacts with the RawData flag. If RawData flag is unset then this flag is ignored.
	// If RawData is `true`, then this flag indicates whether the raw data should be parsed as json.RawMessage
	// and routed via the client's internal logic (`BypassRawDataRouting=false`), or returned to the application code as []byte (`BypassRawDataRouting=true`).
	// If this flag is `true`, it's up to the caller to handle all message types including auth and subscription responses.
	BypassRawDataRouting bool

	// Log is an optional logger. Any logger implementation can be used as long as it
	// implements the basic Logger interface. Omitting this will disable client logging.
	Log Logger
}

func (c *Config) validate() error {
	if c.APIKey == "" {
		return errors.New("API key is required")
	}

	if c.Log == nil {
		c.Log = &nopLogger{}
	}

	return nil
}

// Feed is the data feed (e.g. Delayed, RealTime) which represents the server host.
type Feed string

const (
	Delayed                           Feed = "wss://delayed.polygon.io"
	RealTime                          Feed = "wss://socket.polygon.io"
	Nasdaq                            Feed = "wss://nasdaqfeed.polygon.io"
	PolyFeed                          Feed = "wss://polyfeed.polygon.io"
	PolyFeedPlus                      Feed = "wss://polyfeedplus.polygon.io"
	StarterFeed                       Feed = "wss://starterfeed.polygon.io"
	LaunchpadFeed                     Feed = "wss://launchpad.polygon.io"
	BusinessFeed                      Feed = "wss://business.polygon.io"
	EdgxBusinessFeed                  Feed = "wss://edgx-business.polygon.io"
	DelayedBusinessFeed               Feed = "wss://delayed-business.polygon.io"
	DelayedEdgxBusinessFeed           Feed = "wss://delayed-edgx-business.polygon.io"
	DelayedNasdaqLastSaleBusinessFeed Feed = "wss://delayed-nasdaq-last-sale-business.polygon.io"
	DelayedNasdaqBasicFeed            Feed = "wss://delayed-nasdaq-basic-business.polygon.io"
	DelayedFullMarketBusinessFeed     Feed = "wss://delayed-fullmarket-business.polygon.io"
	FullMarketBusinessFeed            Feed = "wss://fullmarket-business.polygon.io"
	NasdaqLastSaleBusinessFeed        Feed = "wss://nasdaq-last-sale-business.polygon.io"
	NasdaqBasicBusinessFeed           Feed = "wss://nasdaq-basic-business.polygon.io"
)

// Market is the type of market (e.g. Stocks, Crypto) used to connect to the server.
type Market string

const (
	Stocks  Market = "stocks"
	Options Market = "options"
	Forex   Market = "forex"
	Crypto  Market = "crypto"
	Indices Market = "indices"
)

func (m Market) supports(topic Topic) bool {
	switch m {
	case Stocks:
		return topic > stocksMin && topic < stocksMax
	case Options:
		return topic > optionsMin && topic < optionsMax
	case Forex:
		return topic > forexMin && topic < forexMax
	case Crypto:
		return topic > cryptoMin && topic < cryptoMax
	}
	return true // assume user knows what they're doing if they use some unknown market
}

// Topic is the data type used to subscribe and retrieve data from the server.
type Topic uint8

// The launchpad topics should be used for any asset class when connecting to
// the Launchpad feed
const (
	stocksMin              Topic = 10
	StocksSecAggs          Topic = 11
	StocksMinAggs          Topic = 12
	StocksTrades           Topic = 13
	StocksQuotes           Topic = 14
	StocksImbalances       Topic = 15
	StocksLULD             Topic = 16
	StocksLaunchpadMinAggs Topic = 17
	StocksLaunchpadValue   Topic = 18
	stocksMax              Topic = 19

	optionsMin              Topic = 30
	OptionsSecAggs          Topic = 31
	OptionsMinAggs          Topic = 32
	OptionsTrades           Topic = 33
	OptionsQuotes           Topic = 34
	OptionsLaunchpadMinAggs Topic = 35
	OptionsLaunchpadValue   Topic = 36
	optionsMax              Topic = 37

	forexMin              Topic = 50
	ForexSecAggs          Topic = 51
	ForexMinAggs          Topic = 52
	ForexQuotes           Topic = 53
	ForexLaunchpadMinAggs Topic = 54
	ForexLaunchpadValue   Topic = 55
	forexMax              Topic = 56

	cryptoMin              Topic = 70
	CryptoSecAggs          Topic = 71
	CryptoMinAggs          Topic = 72
	CryptoTrades           Topic = 73
	CryptoQuotes           Topic = 74
	CryptoL2Book           Topic = 75
	CryptoLaunchpadMinAggs Topic = 76
	CryptoLaunchpadValue   Topic = 77
	cryptoMax              Topic = 78

	IndexMinAggs Topic = 90
	IndexValue   Topic = 91

	BusinessFairMarketValue Topic = 100
)

func (t Topic) prefix() string {
	switch t {
	case StocksSecAggs:
		return "A"
	case StocksMinAggs:
		return "AM"
	case StocksTrades:
		return "T"
	case StocksQuotes:
		return "Q"
	case StocksImbalances:
		return "NOI"
	case StocksLULD:
		return "LULD"
	case StocksLaunchpadMinAggs:
		return "AM"
	case StocksLaunchpadValue:
		return "LV"
	case OptionsSecAggs:
		return "A"
	case OptionsMinAggs:
		return "AM"
	case OptionsTrades:
		return "T"
	case OptionsQuotes:
		return "Q"
	case OptionsLaunchpadMinAggs:
		return "AM"
	case OptionsLaunchpadValue:
		return "LV"
	case ForexSecAggs:
		return "CAS"
	case ForexMinAggs:
		return "CA"
	case ForexQuotes:
		return "C"
	case ForexLaunchpadMinAggs:
		return "AM"
	case ForexLaunchpadValue:
		return "LV"
	case CryptoSecAggs:
		return "XAS"
	case CryptoMinAggs:
		return "XA"
	case CryptoTrades:
		return "XT"
	case CryptoQuotes:
		return "XQ"
	case CryptoL2Book:
		return "XL2"
	case CryptoLaunchpadMinAggs:
		return "AM"
	case CryptoLaunchpadValue:
		return "LV"
	case IndexMinAggs:
		return "AM"
	case IndexValue:
		return "V"
	case BusinessFairMarketValue:
		return "FMV"
	}
	return ""
}

// Logger is a basic logger interface used for logging within the client.
type Logger interface {
	Debugf(template string, args ...any)
	Infof(template string, args ...any)
	Errorf(template string, args ...any)
}

type nopLogger struct{}

func (l *nopLogger) Debugf(template string, args ...any) {}
func (l *nopLogger) Infof(template string, args ...any)  {}
func (l *nopLogger) Errorf(template string, args ...any) {}
