package models

import "time"

// ListTradesParams is the set of parameters for the ListTrades method.
type ListTradesParams struct {
	// The ticker symbol to get trades for.
	Ticker string `validate:"required" path:"ticker"`

	// Query by timestamp. To query for a specific day instead of a nanosecond timestamp,
	// set it via this pattern: params.WithDay(2006, 1, 2) // January 2, 2006.
	TimestampEQ  *Nanos `query:"timestamp"`
	TimestampLT  *Nanos `query:"timestamp.lt"`
	TimestampLTE *Nanos `query:"timestamp.lte"`
	TimestampGT  *Nanos `query:"timestamp.gt"`
	TimestampGTE *Nanos `query:"timestamp.gte"`

	// Order results based on the sort field.
	Order *Order `query:"order"`

	// Limit the number of results returned, default is 10 and max is 50000.
	Limit *int `query:"limit"`

	// Sort field used for ordering.
	Sort *Sort `query:"sort"`
}

func (p ListTradesParams) WithTimestamp(c Comparator, q Nanos) *ListTradesParams {
	switch c {
	case EQ:
		p.TimestampEQ = &q
	case LT:
		p.TimestampLT = &q
	case LTE:
		p.TimestampLTE = &q
	case GT:
		p.TimestampGT = &q
	case GTE:
		p.TimestampGTE = &q
	}
	return &p
}

func (p ListTradesParams) WithDay(year int, month time.Month, day int) *ListTradesParams {
	d := Nanos(time.Date(year, month, day, 0, 0, 0, 0, time.UTC))
	p.TimestampEQ = &d
	return &p
}

func (p ListTradesParams) WithOrder(q Order) *ListTradesParams {
	p.Order = &q
	return &p
}

func (p ListTradesParams) WithLimit(q int) *ListTradesParams {
	p.Limit = &q
	return &p
}

func (p ListTradesParams) WithSort(q Sort) *ListTradesParams {
	p.Sort = &q
	return &p
}

// ListTradesResponse is the response returned by the ListTrades method.
type ListTradesResponse struct {
	BaseResponse
	Results []Trade `json:"results,omitempty"`
}

// GetLastTradeParams is the set of parameters for GetLastTrade method.
type GetLastTradeParams struct {
	// The ticker symbol of the stock/equity.
	Ticker string `validate:"required" path:"ticker"`
}

// GetLastTradeResponse is the response returned by the GetLastTradeResponse method.
type GetLastTradeResponse struct {
	BaseResponse
	Results LastTrade `json:"results,omitempty"`
}

// GetLastCryptoTradeParams is the set of parameters for the GetLastCryptoTrade method.
type GetLastCryptoTradeParams struct {
	// The "from" symbol of the pair.
	From string `validate:"required" path:"from"`

	// The "to" symbol of the pair.
	To string `validate:"required" path:"to"`
}

// GetLastCryptoTradeResponse is the response returned by the GetLastCryptoTrade method.
type GetLastCryptoTradeResponse struct {
	BaseResponse
	Symbol string      `json:"symbol,omitempty"`
	Last   CryptoTrade `json:"last,omitempty"`
}

// Trade contains trade data for a specified ticker symbol.
type Trade struct {
	Conditions           []int32       `json:"conditions,omitempty"`
	Correction           int           `json:"correction,omitempty"`
	Exchange             int           `json:"exchange,omitempty"`
	ID                   string        `json:"id,omitempty"`
	ParticipantTimestamp Nanos         `json:"participant_timestamp,omitempty"`
	Price                float64       `json:"price,omitempty"`
	SequenceNumber       int64         `json:"sequence_number,omitempty"`
	SipTimestamp         Nanos         `json:"sip_timestamp,omitempty"`
	Size                 float64       `json:"size,omitempty"`
	DecimalSize          DecimalVolume `json:"decimal_size,omitempty"`
	Tape                 int32         `json:"tape,omitempty"`
	TrfID                int           `json:"trf_id,omitempty"`
	TrfTimestamp         Nanos         `json:"trf_timestamp,omitempty"`
}

// LastTrade is the most recent trade for a specified ticker.
type LastTrade struct {
	Ticker               string  `json:"T,omitempty"`
	TRFTimestamp         Nanos   `json:"f,omitempty"`
	SequenceNumber       int64   `json:"q,omitempty"`
	Timestamp            Nanos   `json:"t,omitempty"`
	ParticipantTimestamp Nanos   `json:"y,omitempty"`
	Conditions           []int32 `json:"c,omitempty"`
	Correction           uint32  `json:"e,omitempty"`
	ID                   string  `json:"i,omitempty"`
	Price                float64 `json:"p,omitempty"`
	TRF                  int32   `json:"r,omitempty"`
	Size                 float64 `json:"s,omitempty"`
	Exchange             int32   `json:"x,omitempty"`
	Tape                 int32   `json:"z,omitempty"`
}

// CryptoTrade is a trade for a crypto pair.
type CryptoTrade struct {
	Conditions []int   `json:"conditions,omitempty"`
	Exchange   int     `json:"exchange,omitempty"`
	Price      float64 `json:"price,omitempty"`
	Size       float64 `json:"size,omitempty"`
	Timestamp  Nanos   `json:"timestamp,omitempty"`
}
