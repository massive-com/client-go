package models

// ListAggsParams is the set of parameters for the ListAggs method.
type ListAggsParams struct {
	// The ticker symbol of the stock/equity.
	Ticker string `validate:"required" path:"ticker"`

	// The size of the timespan multiplier.
	Multiplier int `validate:"required" path:"multiplier"`

	// The size of the time window.
	Timespan Timespan `validate:"required" path:"timespan"`

	// The start of the aggregate time window. Either a date with the format YYYY-MM-DD or a millisecond timestamp.
	From Millis `validate:"required" path:"from"`

	// The end of the aggregate time window. Either a date with the format YYYY-MM-DD or a millisecond timestamp.
	To Millis `validate:"required" path:"to"`

	// Whether or not the results are adjusted for splits. By default, results are adjusted. Set this to false to get
	// results that are NOT adjusted for splits.
	Adjusted *bool `query:"adjusted"`

	// Order the results by timestamp. asc will return results in ascending order (oldest at the top), desc will return
	// results in descending order (newest at the top).
	Order *Order `query:"sort"`

	// Limits the number of base aggregates queried to create the aggregate results. Max 50000 and Default 5000. Read
	// more about how limit is used to calculate aggregate results in our article on Aggregate Data API Improvements:
	// https://massive.com/blog/aggs-api-updates/.
	Limit *int `query:"limit"`
}

func (p ListAggsParams) WithAdjusted(q bool) *ListAggsParams {
	p.Adjusted = &q
	return &p
}

func (p ListAggsParams) WithOrder(q Order) *ListAggsParams {
	p.Order = &q
	return &p
}

func (p ListAggsParams) WithLimit(q int) *ListAggsParams {
	p.Limit = &q
	return &p
}

// ListAggsResponse is the response returned by the ListAggs method.
type ListAggsResponse struct {
	BaseResponse
	Ticker       string `json:"ticker,omitempty"`
	QueryCount   int    `json:"queryCount,omitempty"`
	ResultsCount int    `json:"resultsCount,omitempty"`
	Adjusted     bool   `json:"adjusted"`
	Results      []Agg  `json:"results,omitempty"`
}

// GetAggsParams is the set of parameters for the GetAggs method.
type GetAggsParams struct {
	// The ticker symbol of the stock/equity.
	Ticker string `validate:"required" path:"ticker"`

	// The size of the timespan multiplier.
	Multiplier int `validate:"required" path:"multiplier"`

	// The size of the time window.
	Timespan Timespan `validate:"required" path:"timespan"`

	// The start of the aggregate time window. Either a date with the format YYYY-MM-DD or a millisecond timestamp.
	From Millis `validate:"required" path:"from"`

	// The end of the aggregate time window. Either a date with the format YYYY-MM-DD or a millisecond timestamp.
	To Millis `validate:"required" path:"to"`

	// Whether or not the results are adjusted for splits. By default, results are adjusted. Set this to false to get
	// results that are NOT adjusted for splits.
	Adjusted *bool `query:"adjusted"`

	// Order the results by timestamp. asc will return results in ascending order (oldest at the top), desc will return
	// results in descending order (newest at the top).
	Order *Order `query:"sort"`

	// Limits the number of base aggregates queried to create the aggregate results. Max 50000 and Default 5000. Read
	// more about how limit is used to calculate aggregate results in our article on Aggregate Data API Improvements:
	// https://massive.com/blog/aggs-api-updates/.
	Limit *int `query:"limit"`
}

func (p GetAggsParams) WithAdjusted(q bool) *GetAggsParams {
	p.Adjusted = &q
	return &p
}

func (p GetAggsParams) WithOrder(q Order) *GetAggsParams {
	p.Order = &q
	return &p
}

func (p GetAggsParams) WithLimit(q int) *GetAggsParams {
	p.Limit = &q
	return &p
}

// GetAggsResponse is the response returned by the GetAggs method.
type GetAggsResponse struct {
	BaseResponse
	Ticker       string `json:"ticker,omitempty"`
	QueryCount   int    `json:"queryCount,omitempty"`
	ResultsCount int    `json:"resultsCount,omitempty"`
	Adjusted     bool   `json:"adjusted"`
	Results      []Agg  `json:"results,omitempty"`
}

// GetGroupedDailyAggsParams is the set of parameters for the GetGroupedDailyAggs method.
type GetGroupedDailyAggsParams struct {
	// The locale of the market.
	Locale MarketLocale `validate:"required" path:"locale"`

	// The type of market to query.
	MarketType MarketType `validate:"required" path:"marketType"`

	// The beginning date for the aggregate window.
	Date Date `validate:"required" path:"date"`

	// Whether or not the results are adjusted for splits. By default, results are adjusted. Set this to false to get
	// results that are NOT adjusted for splits.
	Adjusted *bool `query:"adjusted"`

	// Include OTC securities in the response. Default is false (don't include OTC securities).
	IncludeOTC *bool `query:"include_otc"`
}

func (p GetGroupedDailyAggsParams) WithAdjusted(q bool) *GetGroupedDailyAggsParams {
	p.Adjusted = &q
	return &p
}

func (p GetGroupedDailyAggsParams) WithIncludeOTC(q bool) *GetGroupedDailyAggsParams {
	p.IncludeOTC = &q
	return &p
}

// GetGroupedDailyAggsResponse is the response returned by the GetGroupedDailyAggs method.
type GetGroupedDailyAggsResponse struct {
	BaseResponse
	Ticker       string `json:"ticker,omitempty"`
	QueryCount   int    `json:"queryCount,omitempty"`
	ResultsCount int    `json:"resultsCount,omitempty"`
	Adjusted     bool   `json:"adjusted"`
	Results      []Agg  `json:"results,omitempty"`
}

// GetDailyOpenCloseAggParams is the set of parameters for the GetDailyOpenCloseAgg method.
type GetDailyOpenCloseAggParams struct {
	// The ticker symbol of the stock/equity.
	Ticker string `validate:"required" path:"ticker"`

	// The date of the requested open/close in the format YYYY-MM-DD.
	Date Date `validate:"required" path:"date"`

	// Whether or not the results are adjusted for splits. By default, results are adjusted. Set this to false to get
	// results that are NOT adjusted for splits.
	Adjusted *bool `query:"adjusted"`
}

func (p GetDailyOpenCloseAggParams) WithAdjusted(q bool) *GetDailyOpenCloseAggParams {
	p.Adjusted = &q
	return &p
}

// GetDailyOpenCloseAggResponse is the response for the GetDailyOpenCloseAgg method.
type GetDailyOpenCloseAggResponse struct {
	BaseResponse
	Symbol        string        `json:"symbol,omitempty"`
	From          string        `json:"from,omitempty"`
	Open          float64       `json:"open,omitempty"`
	High          float64       `json:"high,omitempty"`
	Low           float64       `json:"low,omitempty"`
	Close         float64       `json:"close,omitempty"`
	Volume        float64       `json:"volume,omitempty"`
	DecimalVolume DecimalVolume `json:"decimal_volume,omitempty"`
	AfterHours    float64       `json:"afterHours,omitempty"`
	PreMarket     float64       `json:"preMarket,omitempty"`
	OTC           bool          `json:"otc,omitempty"`
}

// GetPreviousCloseAggParams is the set of parameters for the GetPreviousCloseAgg method.
type GetPreviousCloseAggParams struct {
	// The ticker symbol of the stock/equity.
	Ticker string `validate:"required" path:"ticker"`

	// Whether or not the results are adjusted for splits. By default, results are adjusted. Set this to false to get
	// results that are NOT adjusted for splits.
	Adjusted *bool `query:"adjusted"`
}

func (p GetPreviousCloseAggParams) WithAdjusted(q bool) *GetPreviousCloseAggParams {
	p.Adjusted = &q
	return &p
}

// GetPreviousCloseAggResponse is the response returned by the GetPreviousCloseAgg method.
type GetPreviousCloseAggResponse struct {
	BaseResponse
	Ticker       string `json:"ticker,omitempty"`
	QueryCount   int    `json:"queryCount,omitempty"`
	ResultsCount int    `json:"resultsCount,omitempty"`
	Adjusted     bool   `json:"adjusted"`
	Results      []Agg  `json:"results,omitempty"`
}

// Agg is an aggregation of all the activity on a specified ticker between the start and end timestamps.
type Agg struct {
	Ticker       string  `json:"T,omitempty"`
	Close        float64 `json:"c,omitempty"`
	High         float64 `json:"h,omitempty"`
	Low          float64 `json:"l,omitempty"`
	Transactions int64   `json:"n,omitempty"`
	Open         float64 `json:"o,omitempty"`
	Timestamp    Millis  `json:"t,omitempty"`
	Volume       float64 `json:"v,omitempty"`
	VWAP         float64 `json:"vw,omitempty"`
	OTC          bool    `json:"otc,omitempty"`
}
