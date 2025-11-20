package models

// ListStockFinancialsParams is the set of parameters for the ListStockFinancials method.
type ListStockFinancialsParams struct {
	// Query by company ticker.
	Ticker *string `query:"ticker"`

	// Query by central index key Number (CIK: https://www.sec.gov/edgar/searchedgar/cik.htm).
	CIK *string `query:"cik"`

	// Query by company name.
	CompanyNameFull   *string `query:"company_name"`
	CompanyNameSearch *string `query:"company_name.search"`

	// Query by standard industrial classification (SIC:
	// https://www.sec.gov/corpfin/division-of-corporation-finance-standard-industrial-classification-sic-code-list).
	SIC *string `query:"sic"`

	// Query by the date when the filing with financials data was filed in YYYY-MM-DD format.
	//
	// Best used when querying over date ranges to find financials based on filings that happen in a time period.
	//
	// Examples:
	//
	// To get financials based on filings that have happened after January 1, 2009 use the query param
	// filing_date.gte=2009-01-01.
	//
	// To get financials based on filings that happened in the year 2009 use the query params
	// filing_date.gte=2009-01-01&filing_date.lt=2010-01-01.
	FilingDateEQ  *Date `query:"filing_date"`
	FilingDateLT  *Date `query:"filing_date.lt"`
	FilingDateLTE *Date `query:"filing_date.lte"`
	FilingDateGT  *Date `query:"filing_date.gt"`
	FilingDateGTE *Date `query:"filing_date.gte"`

	// The period of report for the filing with financials data in YYYY-MM-DD format.
	PeriodOfReportDateEQ  *Date `query:"period_of_report_date"`
	PeriodOfReportDateLT  *Date `query:"period_of_report_date.lt"`
	PeriodOfReportDateLTE *Date `query:"period_of_report_date.lte"`
	PeriodOfReportDateGT  *Date `query:"period_of_report_date.gt"`
	PeriodOfReportDateGTE *Date `query:"period_of_report_date.gte"`

	// Query by timeframe. Annual financials originate from 10-K filings, and quarterly financials originate from 10-Q
	// filings. Note: Most companies do not file quarterly reports for Q4 and instead include those financials in their
	// annual report, so some companies my not return quarterly financials for Q4.
	Timeframe *Timeframe `query:"timeframe"`

	// Whether or not to include the xpath and formula attributes for each financial data point. See the xpath and
	// formula response attributes for more info. False by default.
	IncludeSources *bool `query:"include_sources"`

	// Order results based on the sort field.
	Order *Order `query:"order"`

	// Limit the number of results returned, default is 10 and max is 100.
	Limit *int `query:"limit"`

	// Sort field used for ordering.
	Sort *Sort `query:"sort"`
}

func (p ListStockFinancialsParams) WithTicker(q string) *ListStockFinancialsParams {
	p.Ticker = &q
	return &p
}

func (p ListStockFinancialsParams) WithCIK(q string) *ListStockFinancialsParams {
	p.CIK = &q
	return &p
}

func (p ListStockFinancialsParams) WithCompanyName(c NameComparator, q string) *ListStockFinancialsParams {
	switch c {
	case Full:
		p.CompanyNameFull = &q
	case Search:
		p.CompanyNameSearch = &q
	}
	return &p
}

func (p ListStockFinancialsParams) WithSIC(q string) *ListStockFinancialsParams {
	p.SIC = &q
	return &p
}

func (p ListStockFinancialsParams) WithFilingDate(c Comparator, q Date) *ListStockFinancialsParams {
	switch c {
	case EQ:
		p.FilingDateEQ = &q
	case LT:
		p.FilingDateLT = &q
	case LTE:
		p.FilingDateLTE = &q
	case GT:
		p.FilingDateGT = &q
	case GTE:
		p.FilingDateGTE = &q
	}
	return &p
}

func (p ListStockFinancialsParams) WithPeriodOfReportDate(c Comparator, q Date) *ListStockFinancialsParams {
	switch c {
	case EQ:
		p.PeriodOfReportDateEQ = &q
	case LT:
		p.PeriodOfReportDateLT = &q
	case LTE:
		p.PeriodOfReportDateLTE = &q
	case GT:
		p.PeriodOfReportDateGT = &q
	case GTE:
		p.PeriodOfReportDateGTE = &q
	}
	return &p
}

func (p ListStockFinancialsParams) WithTimeframe(q Timeframe) *ListStockFinancialsParams {
	p.Timeframe = &q
	return &p
}

func (p ListStockFinancialsParams) WithIncludeSources(q bool) *ListStockFinancialsParams {
	p.IncludeSources = &q
	return &p
}

func (p ListStockFinancialsParams) WithOrder(q Order) *ListStockFinancialsParams {
	p.Order = &q
	return &p
}

func (p ListStockFinancialsParams) WithLimit(q int) *ListStockFinancialsParams {
	p.Limit = &q
	return &p
}

func (p ListStockFinancialsParams) WithSort(q Sort) *ListStockFinancialsParams {
	p.Sort = &q
	return &p
}

// ListStockFinancialsResponse is the response returned by the ListFinancials method.
type ListStockFinancialsResponse struct {
	BaseResponse
	Results []StockFinancial `json:"results,omitempty"`
}

// StockFinancial contains detailed information on a specified stock financial.
type StockFinancial struct {
	CIK                 string               `json:"cik,omitempty"`
	CompanyName         string               `json:"company_name,omitempty"`
	EndDate             string               `json:"end_date,omitempty"`
	FilingDate          string               `json:"filing_date,omitempty"`
	Financials          map[string]Financial `json:"financials,omitempty"`
	FiscalPeriod        string               `json:"fiscal_period,omitempty"`
	FiscalYear          string               `json:"fiscal_year,omitempty"`
	SourceFilingFileUrl string               `json:"source_filing_file_url,omitempty"`
	SourceFilingUrl     string               `json:"source_filing_url,omitempty"`
	StartDate           string               `json:"start_date,omitempty"`
}

// Financial aliases nested data points of information for a stock financial.
type Financial map[string]struct {
	Formula string  `json:"formula,omitempty"`
	Label   string  `json:"label,omitempty"`
	Order   int32   `json:"order,omitempty"`
	Unit    string  `json:"unit,omitempty"`
	Value   float64 `json:"value,omitempty"`
	Xpath   string  `json:"xpath,omitempty"`
}
