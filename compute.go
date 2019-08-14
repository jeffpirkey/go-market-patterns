package main

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"market-patterns/mal"
	"market-patterns/model"
	"strconv"
	"sync"
)

// Deletes all patterns and series with the given length.  Then, computes
// all ticker periods for the given series length.
func truncAndComputeAllSeries(computeLength int) error {

	err := Repos.PatternRepo.DeleteByLength(computeLength)
	if err != nil {
		return errors.Wrap(err, "problems trunc and computing series")
	}

	err = Repos.SeriesRepo.DeleteByLength(computeLength)
	if err != nil {
		return errors.Wrap(err, "problems trunc and computing series")
	}

	return recomputeAllSeries(computeLength)
}

func recomputeAllSeries(computeLength int) error {
	var trainErrors error

	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, 100)
	tickers := Repos.TickerRepo.FindSymbols()
	for _, symbol := range tickers {
		periods, err := Repos.PeriodRepo.FindBySymbol(symbol, mal.SortAsc)
		if err != nil {
			trainErrors = multierror.Append(trainErrors, err)
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			semaphore <- struct{}{} // Lock
			defer func() {
				<-semaphore // Unlock
			}()

			computeSeries(computeLength, symbol, periods)
		}()
	}

	wg.Wait()

	return trainErrors
}

func computeAllSeries(seriesLength int, dataMap map[*model.Ticker][]*model.Period) error {

	var trainErrors error

	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, 100)
	for ticker, periods := range dataMap {

		wg.Add(1)
		go func() {
			defer wg.Done()
			semaphore <- struct{}{} // Lock
			defer func() {
				<-semaphore // Unlock
			}()

			computeSeries(seriesLength, ticker.Symbol, periods)
		}()
	}

	wg.Wait()

	return trainErrors
}

func computeSeries(computeLength int, symbol string, periods []*model.Period) {

	var patternMap = make(map[string]*model.Pattern)
	for i, period := range periods {

		// Skip until we have enough back periods for the pattern sequencing
		// Must have at least series length + 1 to train
		if i <= computeLength {
			continue
		}

		// Previous pattern name, such as 'UUD' for a pattern of Up -> Up -> Down.
		var patName string
		for x := computeLength; x >= 1; x-- {
			patName += fmt.Sprint(periods[i-x].DailyResult)
		}
		r := model.Calc(periods[i-1].Close, period.Close)

		// Find the pattern and increment the total for the given result, r
		var pattern *model.Pattern
		pattern, found := patternMap[patName]
		if !found {
			pattern = &model.Pattern{Symbol: symbol, Value: patName}
			patternMap[patName] = pattern
		}

		switch r {
		case "U":
			pattern.UpCount++
		case "D":
			pattern.DownCount++
		case "N":
			pattern.NoChangeCount++
		}
		pattern.TotalCount++
	}

	if len(patternMap) > 0 {

		var patterns []*model.Pattern
		for _, v := range patternMap {
			patterns = append(patterns, v)
		}

		// TODO jpirkey validate results
		_, err := Repos.PatternRepo.InsertMany(patterns)
		if err != nil {
			log.WithError(err).Warnf("inserting patterns")
		}

		series := &model.Series{Symbol: symbol, Length: computeLength,
			Name: strconv.Itoa(computeLength) + "-period-series"}
		err = Repos.SeriesRepo.InsertOne(series)
		if err != nil {
			log.WithError(err).Warnf("inserting series")
		}
	} else {
		log.Warnf("[%v] No patterns computing", symbol)
	}
}
