package main

import (
	"fmt"
	"market-patterns/model"
)

func summary(tsym string) {
	ticker := Patterns.Find(tsym)

	for k, v := range ticker.FindAll() {

		fmt.Println(k)
		for k2, v2 := range v.FindAll() {

			fmt.Println(fmt.Sprintf("%v: %v", k2, v2))
			//fmt.Println(fmt.Sprintf("%v avg = %.2f for %v", k, float64(v2)/float64(pattern.TotalCount(k)) * 100, k2))
		}
	}
}

func find50(tsym string) []*model.Period {

	var found = make([]*model.Period, 1)
	tickerPatterns := Patterns.Find(tsym)

	for k, v := range tickerPatterns.FindAll() {

		for k2, v2 := range v.FindAll() {
			c := float64(v2) / float64(v.TotalCount()) * 100
			if c >= 50 {
				fmt.Println(fmt.Sprintf("%v avg = %.2f for %v", k, float64(v2)/float64(v.TotalCount())*100, k2))
			}
		}
	}

	return found
}

func findLastPeriod(tsym string) (model.Period, error) {

	tmp := Periods[tsym]
	if tmp == nil {
		return model.Period{}, fmt.Errorf("unable to find ticker data for %v", tsym)
	}

	return tmp.Last(), nil
}