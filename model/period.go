package model

import "time"

// Defines a single time period
type Period struct {
	Value float64
	Date  time.Time
}

type PeriodTimeSlice []Period

func (p PeriodTimeSlice) Len() int {
	return len(p)
}

func (p PeriodTimeSlice) Less(i, j int) bool {
	return p[i].Date.Before(p[j].Date)
}

func (p PeriodTimeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PeriodTimeSlice) Last() Period {
	return p[len(p)-1]
}
