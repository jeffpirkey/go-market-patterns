package main

import (
	"fmt"
	"market-patterns/model"
)

func train(tsym string) {
	periods := Periods[tsym]
	pattern := Patterns.Find(tsym)

	for i, p := range periods {

		if i < 4 {
			// Skipping first ranges
			continue
		}

		// This is the previous range pattern, such as 'UUD'
		var pName string
		for x := 4; x >= 2; x-- {
			t := model.Calc(periods[i-x].Value, periods[i-(x-1)].Value)
			pName += fmt.Sprint(t)
		}

		r := model.Calc(periods[i-1].Value, p.Value)

		p := pattern.Find(pName)
		p.Inc(r)
	}
}
