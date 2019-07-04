package model

import "time"

// Defines a single time period
type Period struct {
	Date   time.Time `json:"date"`
	Open   float64   `json:"open"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Close  float64   `json:"close"`
	Volume int       `json:"volume"`
	// The result of comparison for previous sequence or
	// Not Defined if the first in a chain of sequences.
	SequenceResult Result            `json:"sequence-result"`
	seriesResults  map[string]Result // Maps each series name trained to the calculated result for this period.
}

// Slice of periods used for sorting and other access methods
type PeriodSlice []*Period

// *********************************************************
//   Period methods
// *********************************************************

func (p *Period) AddSeriesResult(seriesName string, r Result) {
	if p.seriesResults == nil {
		p.seriesResults = make(map[string]Result)
	}
	p.seriesResults[seriesName] = r
}

func (p *Period) FindSeriesResult(seriesName string) Result {
	if p.seriesResults == nil {
		p.seriesResults = make(map[string]Result)
	}
	x, found := p.seriesResults[seriesName]
	if !found {
		x = NotDefined
		p.seriesResults[seriesName] = x
	}
	return x
}

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

func (p PeriodSlice) Last() *Period {
	return p[len(p)-1]
}

// Returns pointers to the items from the end of the slice.
func (p PeriodSlice) LastByRange(l int) []*Period {
	return p[(len(p) - l):]
}
