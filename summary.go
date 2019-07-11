package main

import (
	"fmt"
	"market-patterns/model"
)

func summary(tsym string) {

	ticker := Repos.TickerRepo.FindOne(tsym)
	for seqName, pattern := range ticker.FindAllPatterns() {
		fmt.Println(seqName)
		for result, count := range pattern.FindAll() {
			fmt.Println(fmt.Sprintf("%v: %v", result, count))
			//fmt.Println(fmt.Sprintf("%v avg = %.2f for %v", k, float64(v2)/float64(pattern.TotalCount(k)) * 100, k2))
		}
	}
}

func find50(symbol string) []*model.Period {

	var found = make([]*model.Period, 1)
	ticker := Repos.TickerRepo.FindOne(symbol)
	for seqName, pattern := range ticker.FindAllPatterns() {
		for result, count := range pattern.FindAll() {
			c := float64(count) / float64(pattern.TotalCount()) * 100
			if c >= 50 {
				fmt.Println(fmt.Sprintf("%v avg = %.2f for %v",
					seqName, float64(count)/float64(pattern.TotalCount())*100, result))
			}
		}
	}

	return found
}

func findLastPeriod(symbol string) (*model.Period, error) {

	ticker := Repos.TickerRepo.FindOne(symbol)
	slice := ticker.PeriodSlice()
	return slice.Last(), nil
}
