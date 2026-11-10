package massivews

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

	// ReconnectCallback is a callback that is triggered on automatic reconnects by the websocket client.
	// This can be useful for implementing additional logic around reconnect paths e.g. logging, metrics
	// or managing the connection. The callback function takes as input an error type which will be non-nil
	// if the reconnect attempt has failed and is being retried, and will be nil on reconnect success.
	ReconnectCallback func(error)

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
	Delayed                           Feed = "wss://delayed.massive.com"
	RealTime                          Feed = "wss://socket.massive.com"
	Nasdaq                            Feed = "wss://nasdaqfeed.massive.com"
	PolyFeed                          Feed = "wss://polyfeed.massive.com"
	PolyFeedPlus                      Feed = "wss://polyfeedplus.massive.com"
	StarterFeed                       Feed = "wss://starterfeed.massive.com"
	LaunchpadFeed                     Feed = "wss://launchpad.massive.com"
	BusinessFeed                      Feed = "wss://business.massive.com"
	EdgxBusinessFeed                  Feed = "wss://edgx-business.massive.com"
	IEXBusiness                       Feed = "wss://iex-business.massive.com"
	DelayedBusinessFeed               Feed = "wss://delayed-business.massive.com"
	DelayedEdgxBusinessFeed           Feed = "wss://delayed-edgx-business.massive.com"
	DelayedNasdaqLastSaleBusinessFeed Feed = "wss://delayed-nasdaq-last-sale-business.massive.com"
	DelayedNasdaqBasicFeed            Feed = "wss://delayed-nasdaq-basic-business.massive.com"
	DelayedFullMarketBusinessFeed     Feed = "wss://delayed-fullmarket-business.massive.com"
	FullMarketBusinessFeed            Feed = "wss://fullmarket-business.massive.com"
	NasdaqLastSaleBusinessFeed        Feed = "wss://nasdaq-last-sale-business.massive.com"
	NasdaqBasicBusinessFeed           Feed = "wss://nasdaq-basic-business.massive.com"
)

// Market is the type of market (e.g. Stocks, Crypto) used to connect to the server.
type Market string

const (
	Stocks       Market = "stocks"
	Options      Market = "options"
	Forex        Market = "forex"
	Crypto       Market = "crypto"
	Indices      Market = "indices"
	Futures      Market = "futures"
	FuturesCME   Market = "futures/cme"
	FuturesCBOT  Market = "futures/cbot"
	FuturesNYMEX Market = "futures/nymex"
	FuturesCOMEX Market = "futures/comex"
)

func (m Market) supports(topic Topic) bool {
	// FMV is supported for Stocks, Options, Forex, and Crypto but is not within the range
	// so we need to make sure FMV is suppored if these markets are used.
	isBusinessFairMarketValue := topic == BusinessFairMarketValue

	switch m {
	case Stocks:
		return isBusinessFairMarketValue || (topic > stocksMin && topic < stocksMax)
	case Options:
		return isBusinessFairMarketValue || (topic > optionsMin && topic < optionsMax)
	case Forex:
		return isBusinessFairMarketValue || (topic > forexMin && topic < forexMax)
	case Crypto:
		return isBusinessFairMarketValue || (topic > cryptoMin && topic < cryptoMax)
	case Indices:
		return topic == IndexSecAggs || topic == IndexMinAggs || topic == IndexValue
	case Futures, FuturesCME, FuturesCBOT, FuturesNYMEX, FuturesCOMEX:
		return topic > futuresMin && topic < futuresMax
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
	OptionsAdmin            Topic = 37
	optionsMax              Topic = 38

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

	IndexSecAggs Topic = 90
	IndexMinAggs Topic = 91
	IndexValue   Topic = 92

	BusinessFairMarketValue Topic = 100

	futuresMin    Topic = 110
	FutureSecAggs Topic = 111
	FutureMinAggs Topic = 112
	FutureTrades  Topic = 113
	FutureQuotes  Topic = 114
	futuresMax    Topic = 115
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
	case OptionsAdmin:
		return "AD"
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
	case IndexSecAggs:
		return "A"
	case IndexMinAggs:
		return "AM"
	case IndexValue:
		return "V"
	case BusinessFairMarketValue:
		return "FMV"
	case FutureSecAggs:
		return "A"
	case FutureMinAggs:
		return "AM"
	case FutureTrades:
		return "T"
	case FutureQuotes:
		return "Q"
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
