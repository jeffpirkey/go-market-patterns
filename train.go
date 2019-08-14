package main

import (
	"market-patterns/model"
	"sync"
)

// *********************************************************
//   Train day-to-day functions
// *********************************************************

// Trains the day-to-day results for the given tickers and period arrays.  It is assumed that there are
// at least 2 periods for each ticker.
func trainAllDaily(dataMap map[*model.Ticker][]*model.Period) {

	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, 100)
	for _, periods := range dataMap {
		wg.Add(1)
		go func() {
			defer wg.Done()
			semaphore <- struct{}{} // Lock
			defer func() {
				<-semaphore // Unlock
			}()

			trainDaily(periods)
		}()
	}

	wg.Wait()
}

// Trains a range of periods.  Periods array should be sorted ascending.
func trainDaily(periods []*model.Period) {

	// Train the day-to-day results between
	// two consecutive periods across our period slice
	var prev *model.Period
	for i, period := range periods {

		// Set the first index to prev and skip,
		// as we can't compare it to anything
		if i == 0 {
			period.DailyResult = model.NotDefined
			prev = period
			continue
		}

		// Not sure I like this, but skip a previous period that was nil
		if prev == nil {
			prev = period
			continue
		}

		seqResult := model.Calc(prev.Close, period.Close)
		period.DailyResult = seqResult
		// This period become the previous period
		prev = period
	}
}
