package model

import (
	"time"
)

// Defines a single time period
type Period struct {
	Symbol string    `json:"symbol"`
	Date   time.Time `json:"date"`
	Open   float64   `json:"open"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Close  float64   `json:"close"`
	Volume int       `json:"volume"`
	// The result of comparison for previous period or
	// Not Defined if the first in a chain of periods.
	DailyResult string `json:"daily-result"`
}

// Slice of Periods used for sorting and other access methods
type PeriodSlice []*Period

// *********************************************************
//   PeriodSlice methods
// *********************************************************

func (p PeriodSlice) Len() int {
	return len(p)
}

func (p PeriodSlice) Less(i, j int) bool {
	return p[i].Date.Before(p[j].Date)
}

func (p PeriodSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PeriodSlice) Last() Period {
	return *p[len(p)-1]
}

// Returns pointers to the items from the end of the slice.
func (p PeriodSlice) LastByRange(l int) []*Period {
	var tmp []*Period
	tmp = p[(len(p) - l):]
	return tmp
}
