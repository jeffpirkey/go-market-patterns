package core

import (
	"fmt"
)

const (
	Up         = "U"
	NoChange   = "N"
	Down       = "D"
	NotDefined = "x"
)

// For a given sequence of periods, a pattern records the number of times that the next sequence for the Value
// is Up, Down, or NoChange.
type Pattern struct {
	Symbol        string `json:"symbol"`
	Length        int    `json:"length"`
	Value         string `json:"value"`
	UpCount       int    `json:"upCount"`
	DownCount     int    `json:"downCount"`
	NoChangeCount int    `json:"noChangeCount"`
	TotalCount    int    `json:"totalCount"`
}

// This function calculates a Result using 'prev' as the previous value compared to the current date in 'cur'
func Calc(prev, cur float64) string {
	if prev < cur {
		return Up
	} else if prev > cur {
		return Down
	}
	return NoChange
}

type PatternDensity int

const (
	PatternDensityLow PatternDensity = iota
	PatternDensityMedium
	PatternDensityHigh
)

var (
	patternDensityMap = map[string]PatternDensity{"Low": PatternDensityLow,
		"Medium": PatternDensityMedium, "High": PatternDensityHigh}
	patternDensityNames = [...]string{"Low", "Medium", "High"}
)

// *********************************************************
// Utility functions
// *********************************************************

func (density PatternDensity) String() string {
	return patternDensityNames[density]
}

func PatternDensityFromString(str string) (PatternDensity, error) {
	density, found := patternDensityMap[str]
	if !found {
		return density, fmt.Errorf("invalid pattern density enum '%v'", str)
	}

	return density, nil
}

// Slice of Patterns used for sorting and other access methods
type PatternSlice []*Pattern

// *********************************************************
//   PatternSlice methods
// *********************************************************

func (p PatternSlice) Len() int {
	return len(p)
}

func (p PatternSlice) Less(i, j int) bool {
	return p[i].Value < p[j].Value
}

func (p PatternSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PatternSlice) Last() Pattern {
	return *p[len(p)-1]
}
