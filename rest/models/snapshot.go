package models

import (
	"strings"
)

// GetAllTickersSnapshotParams is the set of parameters for the GetAllTickersSnapshot method.
type GetAllTickersSnapshotParams struct {
	// The locale of the market.
	Locale MarketLocale `validate:"required" path:"locale"`

	// The type of market to query.
	MarketType MarketType `validate:"required" path:"marketType"`

	// A comma separated list of tickers to get snapshots for.
	Tickers *string `query:"tickers"`

	// Include OTC securities in the response. Default is false (don't include OTC securities).
	IncludeOTC *bool `query:"include_otc"`
}

func (p GetAllTickersSnapshotParams) WithTickers(q string) *GetAllTickersSnapshotParams {
	p.Tickers = &q
	return &p
}

func (p GetAllTickersSnapshotParams) WithIncludeOTC(q bool) *GetAllTickersSnapshotParams {
	p.IncludeOTC = &q
	return &p
}

// GetAllTickersSnapshotResponse is the response returned by the GetAllTickersSnapshot method.
type GetAllTickersSnapshotResponse struct {
	BaseResponse
	Tickers []TickerSnapshot `json:"tickers,omitempty"`
}

// GetTickerSnapshotParams is the set of parameters for the GetTickerSnapshot method.
type GetTickerSnapshotParams struct {
	// The locale of the market.
	Locale MarketLocale `validate:"required" path:"locale"`

	// The type of market to query.
	MarketType MarketType `validate:"required" path:"marketType"`

	// The ticker symbol of the stock/equity.
	Ticker string `validate:"required" path:"ticker"`
}

// GetTickerSnapshotResponse is the response returned by the GetTickerSnapshot method.
type GetTickerSnapshotResponse struct {
	BaseResponse
	Snapshot TickerSnapshot `json:"ticker,omitempty"`
}

// GetGainersLosersSnapshotParams is the set of parameters for the GetGainersLosersSnapshot method.
type GetGainersLosersSnapshotParams struct {
	// The locale of the market.
	Locale MarketLocale `validate:"required" path:"locale"`

	// The type of market to query.
	MarketType MarketType `validate:"required" path:"marketType"`

	// The direction of the snapshot results to return.
	Direction Direction `validate:"required" path:"direction"`

	// Include OTC securities in the response. Default is false (don't include OTC securities).
	IncludeOTC *bool `query:"include_otc"`
}

func (p GetGainersLosersSnapshotParams) WithIncludeOTC(q bool) *GetGainersLosersSnapshotParams {
	p.IncludeOTC = &q
	return &p
}

// ListOptionsChainParams is a set of parameters for the ListOptionsChainSnapshot method.
type ListOptionsChainParams struct {
	// The underlying ticker symbol of the option contract.
	UnderlyingAsset string `validate:"required" path:"underlyingAsset"`

	// The strike price of the option contract.
	StrikePrice    *float64 `query:"strike_price"`
	StrikePriceLT  *float64 `query:"strike_price.lt"`
	StrikePriceLTE *float64 `query:"strike_price.lte"`
	StrikePriceGT  *float64 `query:"strike_price.gt"`
	StrikePriceGTE *float64 `query:"strike_price.gte"`

	// The type of contract. Can be ContractCall, ContractPut, or in some rare cases, ContractOther.
	ContractType *ContractType `query:"contract_type"`

	// The contract's expiration date in YYYY-MM-DD format.
	ExpirationDateEQ  *Date `query:"expiration_date"`
	ExpirationDateLT  *Date `query:"expiration_date.lt"`
	ExpirationDateLTE *Date `query:"expiration_date.lte"`
	ExpirationDateGT  *Date `query:"expiration_date.gt"`
	ExpirationDateGTE *Date `query:"expiration_date.gte"`

	// Limit the number of results returned, default is 10 and max is 1000.
	Limit *int `query:"limit"`

	// Sort field used for ordering.
	Sort *Sort `query:"sort"`

	// Order results based on the sort field.
	Order *Order `query:"order"`
}

// WithStrikePrice sets strike price to params. Strike Price is the price at which a put or call option can be exercised.
// comparator options include EQ, LT, LTE, GT, and GTE.
// expirationDate should be in YYYY-MM-DD format
func (o ListOptionsChainParams) WithStrikePrice(comparator Comparator, strikePrice float64) *ListOptionsChainParams {
	switch comparator {
	case EQ:
		o.StrikePrice = &strikePrice
	case LT:
		o.StrikePriceLT = &strikePrice
	case LTE:
		o.StrikePriceLTE = &strikePrice
	case GT:
		o.StrikePriceGT = &strikePrice
	case GTE:
		o.StrikePriceGTE = &strikePrice
	default:
		o.StrikePrice = &strikePrice
	}
	return &o
}

// WithContractType sets contract type to params.
// contractType options include ContractCall and ContractPut.
func (o ListOptionsChainParams) WithContractType(contractType ContractType) *ListOptionsChainParams {
	o.ContractType = &contractType
	return &o
}

// WithLimit sets number of results returned. Limit default is 10. Limit must fall in range of 0-1000.
func (o ListOptionsChainParams) WithLimit(limit int) *ListOptionsChainParams {
	o.Limit = &limit
	return &o
}

// WithExpirationDate sets expiration_date query parameter.
// comparator options include EQ, LT, LTE, GT, and GTE.
// expirationDate should be in YYYY-MM-DD format
func (o ListOptionsChainParams) WithExpirationDate(comparator Comparator, expirationDate Date) *ListOptionsChainParams {
	switch comparator {
	case EQ:
		o.ExpirationDateEQ = &expirationDate
	case LT:
		o.ExpirationDateLT = &expirationDate
	case LTE:
		o.ExpirationDateLTE = &expirationDate
	case GT:
		o.ExpirationDateGT = &expirationDate
	case GTE:
		o.ExpirationDateGTE = &expirationDate
	default:
		o.ExpirationDateEQ = &expirationDate
	}
	return &o
}

// WithOrder sets order of results based on the Sort field.
func (o ListOptionsChainParams) WithOrder(order Order) *ListOptionsChainParams {
	o.Order = &order
	return &o
}

// WithSort sets sort field. Sort expects to receive TickerSymbol, ExpirationDate, or StrikePrice as an argument.
func (o ListOptionsChainParams) WithSort(sort Sort) *ListOptionsChainParams {
	switch sort {
	case TickerSymbol:
		o.Sort = &sort
	case ExpirationDate:
		o.Sort = &sort
	case StrikePrice:
		o.Sort = &sort
	}
	return &o
}

type ListOptionsChainSnapshotResponse struct {
	BaseResponse
	Results []OptionContractSnapshot `json:"results,omitempty"`
}

// GetGainersLosersSnapshotResponse is the response returned by the GetGainersLosersSnapshot method.
type GetGainersLosersSnapshotResponse struct {
	BaseResponse
	Tickers []TickerSnapshot `json:"tickers,omitempty"`
}

// GetOptionContractSnapshotParams is the set of parameters for the GetOptionContractSnapshot method.
type GetOptionContractSnapshotParams struct {
	UnderlyingAsset string `validate:"required" path:"underlyingAsset"`
	OptionContract  string `validate:"required" path:"optionContract"`
}

// GetOptionContractSnapshotResponse is the response returned by the GetOptionContractSnapshot method.
type GetOptionContractSnapshotResponse struct {
	BaseResponse
	Results OptionContractSnapshot `json:"results,omitempty"`
}

// GetCryptoFullBookSnapshotParams is the set of parameters for the GetCryptoFullBookSnapshot method.
type GetCryptoFullBookSnapshotParams struct {
	Ticker string `validate:"required" path:"ticker"`
}

// GetCryptoFullBookSnapshotResponse is the response returned by the GetCryptoFullBookSnapshot method.
type GetCryptoFullBookSnapshotResponse struct {
	BaseResponse
	Data FullBookSnapshot `json:"data,omitempty"`
}

// GetIndicesSnapshotParams is the set of parameters for the GetIndicesSnapshot method.
type GetIndicesSnapshotParams struct {
	// The ticker list to get summaries for
	TickerAnyOf *string `query:"ticker.any_of"`
}

func (p GetIndicesSnapshotParams) WithTickerAnyOf(tickers ...string) *GetIndicesSnapshotParams {
	q := strings.Join(tickers, ",")
	p.TickerAnyOf = &q
	return &p
}

// GetIndicesSnapshotResponse is the response returned by the GetIndicesSnapshot method.
type GetIndicesSnapshotResponse struct {
	BaseResponse
	Results []IndexSnapshot `json:"results,omitempty"`
}

// TickerSnapshot is a collection of data for a ticker including the current minute, day, and previous day's aggregate,
// as well as the last trade and quote.
type TickerSnapshot struct {
	Day              DaySnapshot       `json:"day,omitempty"`
	LastQuote        LastQuoteSnapshot `json:"lastQuote,omitempty"`
	LastTrade        LastTradeSnapshot `json:"lastTrade,omitempty"`
	Minute           MinuteSnapshot    `json:"min,omitempty"`
	PrevDay          DaySnapshot       `json:"prevDay,omitempty"`
	Ticker           string            `json:"ticker,omitempty"`
	TodaysChange     float64           `json:"todaysChange,omitempty"`
	TodaysChangePerc float64           `json:"todaysChangePerc,omitempty"`
	Updated          Nanos             `json:"updated,omitempty"`
	FairMarketValue  float64           `json:"fmv,omitempty"`
}

// DaySnapshot is the most recent day agg for a ticker.
type DaySnapshot struct {
	Close                 float64 `json:"c,omitempty"`
	High                  float64 `json:"h,omitempty"`
	Low                   float64 `json:"l,omitempty"`
	Open                  float64 `json:"o,omitempty"`
	Volume                float64 `json:"v,omitempty"`
	VolumeWeightedAverage float64 `json:"vw,omitempty"`
	OTC                   bool    `json:"otc,omitempty"`
}

// LastQuoteSnapshot is the most recent quote for a ticker.
type LastQuoteSnapshot struct {
	AskPrice  float64 `json:"P,omitempty"`
	BidPrice  float64 `json:"p,omitempty"`
	AskSize   float64 `json:"S,omitempty"`
	BidSize   float64 `json:"s,omitempty"`
	Timestamp Nanos   `json:"t,omitempty"`
}

// LastTradeSnapshot is the most recent trade for a ticker.
type LastTradeSnapshot struct {
	Conditions []int   `json:"c,omitempty"`
	TradeID    string  `json:"i,omitempty"`
	Price      float64 `json:"p,omitempty"`
	Size       float64 `json:"s,omitempty"`
	Timestamp  Nanos   `json:"t,omitempty"`
	ExchangeID int     `json:"x,omitempty"`
}

// MinuteSnapshot is the most recent minute agg for a ticker.
type MinuteSnapshot struct {
	AccumulatedVolume     float64 `json:"av,omitempty"`
	Close                 float64 `json:"c,omitempty"`
	High                  float64 `json:"h,omitempty"`
	Low                   float64 `json:"l,omitempty"`
	Open                  float64 `json:"o,omitempty"`
	Volume                float64 `json:"v,omitempty"`
	VolumeWeightedAverage float64 `json:"vw,omitempty"`
	NumberOfTransactions  float64 `json:"n,omitempty"`
	Timestamp             Millis  `json:"t,omitempty"`
	OTC                   bool    `json:"otc,omitempty"`
}

// OptionContractSnapshot is a collection of data for an option contract ticker including the current day aggregate and
// the most recent quote.
type OptionContractSnapshot struct {
	BreakEvenPrice    float64                         `json:"break_even_price,omitempty"`
	Day               DayOptionContractSnapshot       `json:"day,omitempty"`
	Details           OptionDetails                   `json:"details,omitempty"`
	Greeks            Greeks                          `json:"greeks,omitempty"`
	ImpliedVolatility float64                         `json:"implied_volatility,omitempty"`
	LastQuote         LastQuoteOptionContractSnapshot `json:"last_quote,omitempty"`
	LastTrade         LastTradeOptionContractSnapshot `json:"last_trade,omitempty"`
	OpenInterest      float64                         `json:"open_interest,omitempty"`
	UnderlyingAsset   UnderlyingAsset                 `json:"underlying_asset,omitempty"`
	FairMarketValue   float64                         `json:"fmv,omitempty"`
	FMVLastUpdated    Nanos                           `json:"fmv_last_updated,omitempty"`
}

// IndexSnapshot is a collection of data for an index ticker including the current session information and the most recent value.
type IndexSnapshot struct {
	Value        float64      `json:"value,omitempty"`
	Ticker       string       `json:"ticker,omitempty"`
	Name         string       `json:"name,omitempty"`
	Type         string       `json:"type,omitempty"`
	MarketStatus string       `json:"market_status,omitempty"`
	Session      IndexSession `json:"session,omitempty"`
}

type IndexSession struct {
	Change        float64 `json:"change,omitempty"`
	ChangePercent float64 `json:"change_percent,omitempty"`
	Close         float64 `json:"close,omitempty"`
	High          float64 `json:"high,omitempty"`
	Low           float64 `json:"low,omitempty"`
	Open          float64 `json:"open,omitempty"`
	PreviousClose float64 `json:"previous_close,omitempty"`
}

// DayOptionContractSnapshot contains the most recent day agg for an option contract.
type DayOptionContractSnapshot struct {
	Change        float64 `json:"change,omitempty"`
	ChangePercent float64 `json:"change_percent,omitempty"`
	Close         float64 `json:"close,omitempty"`
	High          float64 `json:"high,omitempty"`
	LastUpdated   Nanos   `json:"last_updated,omitempty"`
	Low           float64 `json:"low,omitempty"`
	Open          float64 `json:"open,omitempty"`
	PreviousClose float64 `json:"previous_close,omitempty"`
	Volume        float64 `json:"volume,omitempty"`
	VWAP          float64 `json:"vwap,omitempty"`
}

// OptionDetails contains more detailed information about an option contract.
type OptionDetails struct {
	ContractType      string  `json:"contract_type,omitempty"`
	ExerciseStyle     string  `json:"exercise_style,omitempty"`
	ExpirationDate    Date    `json:"expiration_date,omitempty"`
	SharesPerContract float64 `json:"shares_per_contract,omitempty"`
	StrikePrice       float64 `json:"strike_price,omitempty"`
	Ticker            string  `json:"ticker,omitempty"`
}

// Greeks contains the delta, gamma, vega, and theta of an option contract.
type Greeks struct {
	Delta float64 `json:"delta,omitempty"`
	Gamma float64 `json:"gamma,omitempty"`
	Theta float64 `json:"theta,omitempty"`
	Vega  float64 `json:"vega,omitempty"`
}

// LastQuoteOptionContractSnapshot contains the most recent quote of an option contract.
type LastQuoteOptionContractSnapshot struct {
	Ask         float64 `json:"ask,omitempty"`
	AskSize     float64 `json:"ask_size,omitempty"`
	Bid         float64 `json:"bid,omitempty"`
	BidSize     float64 `json:"bid_size,omitempty"`
	LastUpdated Nanos   `json:"last_updated,omitempty"`
	Midpoint    float64 `json:"midpoint,omitempty"`
	Timeframe   string  `json:"timeframe,omitempty"`
}

type LastTradeOptionContractSnapshot struct {
	Timestamp  Nanos   `json:"sip_timestamp,omitempty"`
	Conditions []int32 `json:"conditions,omitempty"`
	Price      float64 `json:"price,omitempty"`
	Size       float64 `json:"size,omitempty"`
	Exchange   int32   `json:"exchange,omitempty"`
	Timeframe  string  `json:"timeframe,omitempty"`
}

// UnderlyingAsset contains information on the underlying stock for this options contract.
type UnderlyingAsset struct {
	ChangeToBreakEven float64 `json:"change_to_break_even,omitempty"`
	LastUpdated       int64   `json:"last_updated,omitempty"`
	Price             float64 `json:"price,omitempty"`
	Value             float64 `json:"value,omitempty"`
	Ticker            string  `json:"ticker,omitempty"`
	Timeframe         string  `json:"timeframe,omitempty"`
}

// FullBookSnapshot is the level 2 book of a single crypto ticker.
type FullBookSnapshot struct {
	AskCount float64          `json:"askCount,omitempty"`
	Asks     []OrderBookQuote `json:"asks,omitempty"`
	BidCount float64          `json:"bidCount,omitempty"`
	Bids     []OrderBookQuote `json:"bids,omitempty"`
	Spread   float64          `json:"spread,omitempty"`
	Ticker   string           `json:"ticker,omitempty"`
	Updated  Nanos            `json:"updated,omitempty"`
}

// OrderBookQuote contains quote information for a crypto ticker.
type OrderBookQuote struct {
	Price            float64            `json:"p,omitempty"`
	ExchangeToShares map[string]float64 `json:"x,omitempty"`
}

// ListUniversalSnapshotsParams is a set of parameters for the ListUniversalSnapshots method.
type ListUniversalSnapshotsParams struct {
	TickerAnyOf *string `query:"ticker.any_of"`
	Ticker      *string `query:"ticker"`

	TickerLT  *string `query:"ticker.lt"`
	TickerLTE *string `query:"ticker.lte"`
	TickerGT  *string `query:"ticker.gt"`
	TickerGTE *string `query:"ticker.gte"`

	Type *string `query:"type"`
}

// WithTickerAnyOf sets the ticker.any_of query param.
func (p ListUniversalSnapshotsParams) WithTickerAnyOf(q string) *ListUniversalSnapshotsParams {
	p.TickerAnyOf = &q
	return &p
}

// WithTicker sets the ticker equality query param.
func (p ListUniversalSnapshotsParams) WithTicker(q string) *ListUniversalSnapshotsParams {
	p.Ticker = &q
	return &p
}

// WithType sets the type query param.
func (p ListUniversalSnapshotsParams) WithType(q string) *ListUniversalSnapshotsParams {
	p.Type = &q
	return &p
}

// WithTickersByComparison sets the ticker inequality query params.
// Comparator options include EQ, LT, LTE, GT, and GTE.
func (p ListUniversalSnapshotsParams) WithTickersByComparison(c Comparator, q string) *ListUniversalSnapshotsParams {
	switch c {
	case LT:
		p.TickerLT = &q
	case LTE:
		p.TickerLTE = &q
	case GT:
		p.TickerGT = &q
	case GTE:
		p.TickerGTE = &q
	}
	return &p
}

// ListUniversalSnapshotsResponse is the response returned by the ListUniversalSnapshots method.
type ListUniversalSnapshotsResponse struct {
	BaseResponse
	Results []SnapshotResponseModel `json:"results,omitempty"`
}

// SnapshotResponseModel contains all the information that might come back in a SnapshotResponse.
type SnapshotResponseModel struct {
	Name              string             `json:"name,omitempty"`
	MarketStatus      string             `json:"market_status,omitempty"`
	Ticker            string             `json:"ticker,omitempty"`
	Type              string             `json:"type,omitempty"`
	LastMinute        SnapshotLastMinute `json:"last_minute,omitempty"`
	LastQuote         SnapshotLastQuote  `json:"last_quote,omitempty"`
	LastTrade         SnapshotLastTrade  `json:"last_trade,omitempty"`
	Session           Session            `json:"session,omitempty"`
	BreakEvenPrice    float64            `json:"break_even_price,omitempty"`
	Details           Details            `json:"details,omitempty"`
	Greeks            Greeks             `json:"greeks,omitempty"`
	ImpliedVolatility float64            `json:"implied_volatility,omitempty"`
	OpenInterest      float64            `json:"open_interest,omitempty"`
	UnderlyingAsset   UnderlyingAsset    `json:"underlying_asset,omitempty"`
	Value             float64            `json:"value,omitempty"`
	LastUpdated       int64              `json:"last_updated,omitempty"`
	Timeframe         string             `json:"timeframe,omitempty"`
	FairMarketValue   float64            `json:"fmv,omitempty"`
	FMVLastUpdated    Nanos              `json:"fmv_last_updated,omitempty"`
	Error             string             `json:"error"`
	Message           string             `json:"message"`
}

// SnapshotLastMinute contains all the information that might come back in the last_minute attribute of a SnapshotResponseModel.
type SnapshotLastMinute struct {
	Close        float64 `json:"close,omitempty"`
	High         float64 `json:"high,omitempty"`
	Low          float64 `json:"low,omitempty"`
	Transactions int64   `json:"transactions,omitempty"`
	Open         float64 `json:"open,omitempty"`
	Volume       float64 `json:"volume,omitempty"`
	//TODO: check if this need a decimal_volume
	VWAP        float64 `json:"vwap,omitempty"`
	LastUpdated int64   `json:"last_updated,omitempty"`
}

// SnapshotLastQuote contains all the information that might come back in the last_quote attribute of a SnapshotResponse.
type SnapshotLastQuote struct {
	Ask         float64 `json:"ask,omitempty"`
	AskSize     float64 `json:"ask_size,omitempty"`
	Bid         float64 `json:"bid,omitempty"`
	BidSize     float64 `json:"bid_size,omitempty"`
	LastUpdated int64   `json:"last_updated,omitempty"`
	Midpoint    float64 `json:"midpoint,omitempty"`
	Timeframe   string  `json:"timeframe,omitempty"`
	Exchange    int64   `json:"exchange,omitempty"`
}

// SnapshotLastTrade contains all the information that might come back in the last_trade attribute of a SnapshotResponse.
type SnapshotLastTrade struct {
	Timestamp            int64         `json:"sip_timestamp,omitempty"`
	ParticipantTimestamp int64         `json:"participant_timestamp,omitempty"`
	Conditions           []int32       `json:"conditions,omitempty"`
	Price                float64       `json:"price,omitempty"`
	Size                 uint32        `json:"size,omitempty"`
	DecimalSize          DecimalVolume `json:"decimal_size,omitempty"`
	Exchange             int32         `json:"exchange,omitempty"`
	Timeframe            string        `json:"timeframe,omitempty"`
	ID                   string        `json:"id,omitempty"`
	LastUpdated          int64         `json:"last_updated,omitempty"`
}

// Session contains all the information that might come back in the Session attribute of a SnapshotResponseModel or SummaryResult.
type Session struct {
	Change                      float64 `json:"change,omitempty"`
	ChangePercent               float64 `json:"change_percent,omitempty"`
	EarlyTradingChange          float64 `json:"early_trading_change,omitempty"`
	EarlyTradingChangePercent   float64 `json:"early_trading_change_percent,omitempty"`
	RegularTradingChange        float64 `json:"regular_trading_change,omitempty"`
	RegularTradingChangePercent float64 `json:"regular_trading_change_percent,omitempty"`
	LateTradingChange           float64 `json:"late_trading_change,omitempty"`
	LateTradingChangePercent    float64 `json:"late_trading_change_percent,omitempty"`
	Close                       float64 `json:"close,omitempty"`
	High                        float64 `json:"high,omitempty"`
	Low                         float64 `json:"low,omitempty"`
	Open                        float64 `json:"open,omitempty"`
	PreviousClose               float64 `json:"previous_close,omitempty"`
	Volume                      float64 `json:"volume,omitempty"`
	//TODO: check if this need a decimal volume
	Price       float64 `json:"price,omitempty"`
	LastUpdated int64   `json:"last_updated,omitempty"`
	VWAP        float64 `json:"vwap,omitempty"`
}

// Details contains all the information that might come back in the details attribute of a SnapshotResponse.
type Details struct {
	ContractType      string  `json:"contract_type,omitempty"`
	ExerciseStyle     string  `json:"exercise_style,omitempty"`
	ExpirationDate    string  `json:"expiration_date,omitempty"`
	SharesPerContract float64 `json:"shares_per_contract,omitempty"`
	StrikePrice       float64 `json:"strike_price,omitempty"`
}
