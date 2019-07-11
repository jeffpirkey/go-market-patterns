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
	Symbol   string
	Series   map[string]*Series
	Patterns map[string]*Pattern // Map of Patterns by result sequence names, such as 'UUD' for Up/Up/Down
	Periods  map[string]*Period  // Map of Periods by date string
}

// *********************************************************
//   Series methods for Ticker objects
// *********************************************************

func (t *Ticker) AddSeries(name, desc string, len int) *Series {
	newS := &Series{desc, len}
	t.Series[name] = newS
	return newS
}

func (t *Ticker) FindAllSeries() map[string]*Series {
	return t.Series
}

func (t *Ticker) FindSeries(name string) *Series {
	x, found := t.Series[name]
	if !found {
		// Create an empty Series for the given Series name
		x = &Series{}
		t.Series[name] = x
	}
	return x
}

// *********************************************************
//   Pattern methods for Ticker objects
// *********************************************************

func (t *Ticker) FindAllPatterns() map[string]*Pattern {
	return t.Patterns
}

// Finds a pattern with the given name where the name is a combination of Results, such as 'UDD' representing
// the pattern of Up/Down/Down.
func (t *Ticker) FindPattern(patName string) *Pattern {
	x, found := t.Patterns[patName]
	if !found {
		// Create an empty Pattern for the given pattern name
		x = &Pattern{}
		t.Patterns[patName] = x
	}
	return x
}

func (t *Ticker) StartsWithPattern(s string) []*Pattern {
	var p []*Pattern
	for k, v := range t.Patterns {
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
	t.Periods[p.Date.String()] = p
}

func (t *Ticker) FindAllPeriods() map[string]*Period {
	return t.Periods
}

func (t *Ticker) FindPeriod(v time.Time) *Period {
	x, found := t.Periods[v.String()]
	if !found {
		x = &Period{Date: v}
		t.Periods[v.String()] = x
	}
	return x
}

// This function returns a descending sorted slice of the Periods
func (t *Ticker) PeriodSlice() PeriodSlice {
	var slice PeriodSlice
	for _, v := range t.Periods {
		slice = append(slice, v)
	}
	sort.Sort(slice)
	return slice
}
