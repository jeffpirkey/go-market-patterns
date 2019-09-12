package core

type Ticker struct {
	Symbol  string `json:"symbol"`
	Company string `json:"company"`
}

// Slice of Tickers used for sorting and other access methods
type TickerSlice []*Ticker

// *********************************************************
//   PeriodSlice methods
// *********************************************************

func (p TickerSlice) Len() int {
	return len(p)
}

func (p TickerSlice) Less(i, j int) bool {
	return p[i].Symbol < p[j].Symbol
}

func (p TickerSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p TickerSlice) Last() *Ticker {
	return p[len(p)-1]
}

// Returns pointers to the items from the end of the slice.
func (p TickerSlice) LastByRange(l int) []*Ticker {
	return p[(len(p) - l):]
}
