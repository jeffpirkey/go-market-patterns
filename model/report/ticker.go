package report

type TickerSymbolNames struct {
	Names []string `json:"names"`
}

type Tickers struct {
	Tickers *TickerSymbolCompanySlice `json:"tickers"`
}

type TickerSymbolCompany struct {
	Symbol  string `json:"symbol"`
	Company string `json:"company"`
}

// Slice of Ticker Symbols and Company Names used for sorting and other access methods
type TickerSymbolCompanySlice []*TickerSymbolCompany

// *********************************************************
//   PeriodSlice methods
// *********************************************************

func (p TickerSymbolCompanySlice) Len() int {
	return len(p)
}

func (p TickerSymbolCompanySlice) Less(i, j int) bool {
	return p[i].Symbol < p[j].Symbol
}

func (p TickerSymbolCompanySlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p TickerSymbolCompanySlice) Last() *TickerSymbolCompany {
	return p[len(p)-1]
}

// Returns pointers to the items from the end of the slice.
func (p TickerSymbolCompanySlice) LastByRange(l int) []*TickerSymbolCompany {
	return p[(len(p) - l):]
}
