package model

import (
	"sort"
	"strings"
	"time"
)

type Series struct {
	Description  string
	SeriesLength int
}

type Ticker struct {
	series   map[string]*Series
	patterns map[string]*Pattern   // Map of Patterns by result sequence names, such as 'UUD' for Up/Up/Down
	periods  map[time.Time]*Period // Map of Periods by date
}

type Tickers struct {
	tickers map[string]*Ticker
}

func NewTickers() Tickers {
	return Tickers{make(map[string]*Ticker)}
}

// *********************************************************
//   Tickers methods
// *********************************************************

func (t *Tickers) FindNames() []string {

	var names []string
	for name, _ := range t.tickers {
		names = append(names, name)
	}
	sort.Strings(names)

	return names
}

func (t *Tickers) FindAll() map[string]*Ticker {
	return t.tickers
}

func (t *Tickers) Find(tickerSym string) *Ticker {
	x, found := t.tickers[tickerSym]
	if !found {
		// Create an empty Ticker for the given symbol
		x = &Ticker{make(map[string]*Series), make(map[string]*Pattern),
			make(map[time.Time]*Period)}
		t.tickers[tickerSym] = x
	}
	return x
}

// *********************************************************
//   Series methods for Ticker objects
// *********************************************************

func (t *Ticker) AddSeries(name, desc string, len int) *Series {
	newS := &Series{desc, len}
	t.series[name] = newS
	return newS
}

func (t *Ticker) FindAllSeries() map[string]*Series {
	return t.series
}

func (t *Ticker) FindSeries(name string) *Series {
	x, found := t.series[name]
	if !found {
		// Create an empty series for the given series name
		x = &Series{}
		t.series[name] = x
	}
	return x
}

// *********************************************************
//   Pattern methods for Ticker objects
// *********************************************************

func (t *Ticker) FindAllPatterns() map[string]*Pattern {
	return t.patterns
}

// Finds a pattern with the given name where the name is a combination of Results, such as 'UDD' representing
// the pattern of Up/Down/Down.
func (t *Ticker) FindPattern(patName string) *Pattern {
	x, found := t.patterns[patName]
	if !found {
		// Create an empty Pattern for the given pattern name
		x = &Pattern{}
		t.patterns[patName] = x
	}
	return x
}

func (t *Ticker) StartsWithPattern(s string) []*Pattern {
	var p []*Pattern
	for k, v := range t.patterns {
		if strings.HasPrefix(k, s) {
			p = append(p, v)
		}
	}
	return p
}

// *********************************************************
//   Period methods for Ticker objects
// *********************************************************

// Adds a period using the given period's date as the timestamp used when adding the to period map. Replaces
// any value that may have already been there.
func (t *Ticker) AddPeriod(p *Period) {
	t.periods[p.Date] = p
}

func (t *Ticker) FindAllPeriods() map[time.Time]*Period {
	return t.periods
}

func (t *Ticker) FindPeriod(v time.Time) *Period {
	x, found := t.periods[v]
	if !found {
		x = &Period{Date: v}
		t.periods[v] = x
	}
	return x
}

// This function returns a descending sorted slice of the periods
func (t *Ticker) PeriodSlice() PeriodSlice {
	var slice PeriodSlice
	for _, v := range t.periods {
		slice = append(slice, v)
	}
	sort.Sort(slice)
	return slice
}
