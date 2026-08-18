package models

// Action is the set of recognized actions used in control messages.
type Action string

const (
	Auth        Action = "auth"
	Subscribe   Action = "subscribe"
	Unsubscribe Action = "unsubscribe"
)

// EventType is the type of message received. It should be present in every message sent by the server.
type EventType struct {
	EventType string `json:"ev,omitempty"`
}

// ControlMessage is a message to signal status and control events to and from the server.
type ControlMessage struct {
	EventType
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Action  Action `json:"action,omitempty"`
	Params  string `json:"params,omitempty"`
}

// EquityAgg is an aggregate for either stock tickers or option contracts.
type EquityAgg struct {
	// The event type.
	EventType

	// The ticker symbol for the given stock.
	Symbol string `json:"sym,omitempty"`

	// The tick volume.
	Volume float64 `json:"v,omitempty"`

	// The tick volume including fractional shares, represented as a decimal string. Stocks only. Kept as a string so
	// that the precision the server sends is preserved; Volume truncates any fractional part.
	DecimalVolume string `json:"dv,omitempty"`

	// Today's accumulated volume.
	AccumulatedVolume float64 `json:"av,omitempty"`

	// Today's accumulated volume including fractional shares, represented as a decimal string. Stocks only. See
	// DecimalVolume for why this is a string.
	DecimalAccumulatedVolume string `json:"dav,omitempty"`

	// Today's official opening price.
	OfficialOpenPrice float64 `json:"op,omitempty"`

	// The tick's volume weighted average price.
	VWAP float64 `json:"vw,omitempty"`

	// The opening tick price for this aggregate window.
	Open float64 `json:"o,omitempty"`

	// The closing tick price for this aggregate window.
	Close float64 `json:"c,omitempty"`

	// The highest tick price for this aggregate window.
	High float64 `json:"h,omitempty"`

	// The lowest tick price for this aggregate window.
	Low float64 `json:"l,omitempty"`

	// Today's volume weighted average price.
	AggregateVWAP float64 `json:"a,omitempty"`

	// The average trade size for this aggregate window.
	AverageSize float64 `json:"z,omitempty"`

	// The timestamp of the starting tick for this aggregate window in Unix Milliseconds.
	StartTimestamp int64 `json:"s,omitempty"`

	// The timestamp of the ending tick for this aggregate window in Unix Milliseconds.
	EndTimestamp int64 `json:"e,omitempty"`

	// Whether or not this aggregate is for an OTC ticker. This field will be left off if false.
	OTC bool `json:"otc,omitempty"`
}

// CurrencyAgg is an aggregate for either forex currency pairs or crypto pairs.
type CurrencyAgg struct {
	// The event type.
	EventType

	// The currency pair.
	Pair string `json:"pair,omitempty"`

	// The open price for this aggregate window.
	Open float64 `json:"o,omitempty"`

	// The close price for this aggregate window.
	Close float64 `json:"c,omitempty"`

	// The high price for this aggregate window.
	High float64 `json:"h,omitempty"`

	// The low price for this aggregate window.
	Low float64 `json:"l,omitempty"`

	// The volume of trades during this aggregate window.
	Volume float64 `json:"v,omitempty"`

	// The volume weighted average price.
	VWAP float64 `json:"vw,omitempty"`

	// The start time for this aggregate window in Unix Milliseconds.
	StartTimestamp int64 `json:"s,omitempty"`

	// The end time for this aggregate window in Unix Milliseconds.
	EndTimestamp int64 `json:"e,omitempty"`

	// The average trade size for this aggregate window.
	AVGTradeSize int32 `json:"z,omitempty"`
}

// EquityTrade is trade data for either stock tickers or option contracts.
type EquityTrade struct {
	// The event type.
	EventType

	// The ticker symbol for the given stock.
	Symbol string `json:"sym,omitempty"`

	// The exchange ID.
	Exchange int32 `json:"x,omitempty"`

	// The trade ID.
	ID string `json:"i,omitempty"`

	// The tape. (1 = NYSE, 2 = AMEX, 3 = Nasdaq).
	Tape int32 `json:"z,omitempty"`

	// The price.
	Price float64 `json:"p,omitempty"`

	// The trade size.
	Size int64 `json:"s,omitempty"`

	// The trade size including fractional shares, represented as a decimal string. Stocks only. Prefer this over Size
	// for stocks: Size truncates, so a trade of less than one share is reported as 0.
	DecimalSize string `json:"ds,omitempty"`

	// The trade conditions.
	Conditions []int32 `json:"c,omitempty"`

	// The Timestamp in Unix MS.
	Timestamp int64 `json:"t,omitempty"`

	// The time the trade occurred at the exchange or TRF, in Unix MS. This is always less than or equal to Timestamp,
	// which is stamped later by the SIP.
	ParticipantTimestamp int64 `json:"pt,omitempty"`

	// The sequence number represents the sequence in which message events happened. These are increasing and unique per
	// ticker symbol, but will not always be sequential (e.g., 1, 2, 6, 9, 10, 11).
	SequenceNumber int64 `json:"q,omitempty"`

	// The ID for the Trade Reporting Facility where the trade took place.
	TradeReportingFacilityID int64 `json:"trfi,omitempty"`

	// The TRF (Trade Reporting Facility) Timestamp in Unix MS.
	// This is the timestamp of when the trade reporting facility received this trade.
	TradeReportingFacilityTimestamp int64 `json:"trft,omitempty"`
}

// CryptoTrade is a trade for a crypto pair.
type CryptoTrade struct {
	// The event type.
	EventType

	// The crypto pair.
	Pair string `json:"pair,omitempty"`

	// The crypto exchange ID.
	Exchange int32 `json:"x,omitempty"`

	// The ID of the trade (optional).
	ID string `json:"i,omitempty"`

	// The price.
	Price float64 `json:"p,omitempty"`

	// The size.
	Size float64 `json:"s,omitempty"`

	// The conditions. 0 (or empty array): empty 1: sellside 2: buyside
	Conditions []int32 `json:"c,omitempty"`

	// The Timestamp in Unix MS.
	Timestamp int64 `json:"t,omitempty"`

	// The timestamp that the tick was received by Massive.
	ReceivedTimestamp int64 `json:"r,omitempty"`
}

// EquityQuote is a quote for either stock tickers or option contracts.
type EquityQuote struct {
	// The event type.
	EventType

	// The ticker symbol for the given stock.
	Symbol string `json:"sym,omitempty"`

	// The bid exchange ID.
	BidExchangeID int32 `json:"bx,omitempty"`

	// The bid price.
	BidPrice float64 `json:"bp,omitempty"`

	// The bid size. This represents the number of round lot orders at the given bid price. The normal round lot size is
	// 100 shares. A bid size of 2 means there are 200 shares for purchase at the given bid price.
	BidSize int32 `json:"bs,omitempty"`

	// The ask exchange ID.
	AskExchangeID int32 `json:"ax,omitempty"`

	// The ask price.
	AskPrice float64 `json:"ap,omitempty"`

	// The ask size. This represents the number of round lot orders at the given ask price. The normal round lot size is
	// 100 shares. An ask size of 2 means there are 200 shares available to purchase at the given ask price.
	AskSize int32 `json:"as,omitempty"`

	// The condition.
	Condition int32 `json:"c,omitempty"`

	// The indicators. For more information, see our glossary: https://massive.com/glossary/us/stocks/conditions-indicators.
	Indicators []int32 `json:"i,omitempty"`

	// The Timestamp in Unix MS.
	Timestamp int64 `json:"t,omitempty"`

	// The tape. (1 = NYSE, 2 = AMEX, 3 = Nasdaq).
	Tape int32 `json:"z,omitempty"`

	// The sequence number represents the sequence in which message events happened. These are increasing and unique per
	// ticker symbol, but will not always be sequential (e.g., 1, 2, 6, 9, 10, 11).
	SequenceNumber int64 `json:"q,omitempty"`
}

// ForexQuote is a quote for a forex currency pair.
type ForexQuote struct {
	// The event type.
	EventType

	// The current pair.
	Pair string `json:"p,omitempty"`

	// The exchange ID.
	ExchangeID int32 `json:"x,omitempty"`

	// The ask price.
	AskPrice float64 `json:"a,omitempty"`

	// The bid price.
	BidPrice float64 `json:"b,omitempty"`

	// The Timestamp in Unix MS.
	Timestamp int64 `json:"t,omitempty"`
}

// CryptoQuote is a quote for a crypto pair.
type CryptoQuote struct {
	// The event type.
	EventType

	// The crypto pair.
	Pair string `json:"pair,omitempty"`

	// The bid price.
	BidPrice float64 `json:"bp,omitempty"`

	// The bid size.
	BidSize float64 `json:"bs,omitempty"`

	// The ask price.
	AskPrice float64 `json:"ap,omitempty"`

	// The ask size.
	AskSize float64 `json:"as,omitempty"`

	// The Timestamp in Unix MS.
	Timestamp int64 `json:"t,omitempty"`

	// The crypto exchange ID.
	ExchangeID int32 `json:"x,omitempty"`

	// The timestamp that the tick was received by Massive.
	ReceivedTimestamp int64 `json:"r,omitempty"`
}

// Imbalance is an imbalance event for a given stock ticker symbol.
type Imbalance struct {
	// The event type.
	EventType

	// The ticker symbol for the given stock.
	Symbol string `json:"T,omitempty"`

	// The Timestamp in Unix MS.
	Timestamp int64 `json:"t,omitempty"`

	// The time that the auction is planned to take place in the format (hour x 100) + minutes in Eastern Standard Time,
	// for example 930 would be 9:30 am EST, and 1600 would be 4:00 pm EST.
	AuctionTime int32 `json:"at,omitempty"`

	// The auction type. O - Early Opening Auction (non-NYSE only) M - Core Opening Auction H - Reopening Auction (Halt
	// Resume) C - Closing Auction P - Extreme Closing Imbalance (NYSE only) R - Regulatory Closing Imbalance (NYSE
	// only)
	AuctionType string `json:"a,omitempty"`

	// The symbol sequence.
	SymbolSequence int32 `json:"i,omitempty"`

	// The exchange ID.
	ExchangeID int32 `json:"x,omitempty"`

	// The imbalance quantity.
	ImbalanceQuantity int32 `json:"o,omitempty"`

	// The paired quantity.
	PairedQuantity int32 `json:"p,omitempty"`

	// The book clearing price.
	BookClearingPrice float64 `json:"b,omitempty"`
}

// LimitUpLimitDown is a LULD event for a given stock ticker symbol.
type LimitUpLimitDown struct {
	// The event type.
	EventType

	// The ticker symbol for the given stock.
	Symbol string `json:"T,omitempty"`

	// The high price.
	HighPrice float64 `json:"h,omitempty"`

	// The low price.
	LowPrice float64 `json:"l,omitempty"`

	// The Indicators.
	Indicators []int32 `json:"i,omitempty"`

	// The tape. (1 = NYSE, 2 = AMEX, 3 = Nasdaq).
	Tape int32 `json:"z,omitempty"`

	// The Timestamp in Unix MS.
	Timestamp int64 `json:"t,omitempty"`

	// The sequence number represents the sequence in which message events happened. These are increasing and unique per ticker symbol, but will not always be sequential (e.g., 1, 2, 6, 9, 10, 11).
	SequenceNumber int64 `json:"q,omitempty"`
}

// Level2Book is level 2 book data for a given crypto pair.
type Level2Book struct {
	// The event type.
	EventType

	// The crypto pair.
	Pair string `json:"pair,omitempty"`

	// An array of bid prices with a maximum depth of 100.
	BidPrices [][]float64 `json:"b,omitempty"`

	// An array of ask prices with a maximum depth of 100.
	AskPrices [][]float64 `json:"a,omitempty"`

	// The Timestamp in Unix MS.
	Timestamp int64 `json:"t,omitempty"`

	// The crypto exchange ID.
	ExchangeID int32 `json:"x,omitempty"`

	// The timestamp that the tick was received by Massive.
	ReceivedTimestamp int64 `json:"r,omitempty"`
}

// IndexValue is value data for either indices.
type IndexValue struct {
	// The event type.
	EventType

	// The value.
	Value float64 `json:"val"`

	// The ticker symbol for the given index.
	Ticker string `json:"T"`

	// The Timestamp in Unix MS.
	Timestamp int64 `json:"t,omitempty"`
}

type LaunchpadValue struct {
	// The event type.
	EventType

	// The value.
	Value float64 `json:"val"`

	// The ticker symbol for the given security.
	Ticker string `json:"sym"`

	// The Timestamp in nanoseconds.
	Timestamp int64 `json:"t,omitempty"`
}

type FairMarketValue struct {
	// The event type.
	EventType

	// The value.
	FMV float64 `json:"fmv"`

	// The ticker symbol for the given security.
	Ticker string `json:"sym"`

	// The Timestamp in nanoseconds.
	Timestamp int64 `json:"t,omitempty"`
}

// FuturesTrade represents a futures trade event.
type FuturesTrade struct {
	EventType
	Symbol         string  `json:"sym,omitempty"`
	Price          float64 `json:"p,omitempty"`
	Size           int64   `json:"s,omitempty"`
	Timestamp      int64   `json:"t,omitempty"`
	SequenceNumber int64   `json:"q,omitempty"`
}

// FuturesQuote represents a futures quote event.
type FuturesQuote struct {
	EventType
	Symbol       string  `json:"sym,omitempty"`
	BidPrice     float64 `json:"bp,omitempty"`
	BidSize      int64   `json:"bs,omitempty"`
	BidTimestamp int64   `json:"bt,omitempty"`
	AskPrice     float64 `json:"ap,omitempty"`
	AskSize      int64   `json:"as,omitempty"`
	AskTimestamp int64   `json:"at,omitempty"`
	Timestamp    int64   `json:"t,omitempty"`
}

// FuturesAggregate represents an aggregate event (e.g., second or minute) for a futures contract.
type FuturesAggregate struct {
	EventType      string  `json:"ev,omitempty"`
	Symbol         string  `json:"sym,omitempty"`
	Volume         float64 `json:"v,omitempty"`
	TotalValue     float64 `json:"dv,omitempty"`
	Open           float64 `json:"o,omitempty"`
	Close          float64 `json:"c,omitempty"`
	High           float64 `json:"h,omitempty"`
	Low            float64 `json:"l,omitempty"`
	Transactions   int64   `json:"n,omitempty"`
	StartTimestamp int64   `json:"s,omitempty"`
	EndTimestamp   int64   `json:"e,omitempty"`
}
