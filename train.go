package main

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"market-patterns/model"
	"strconv"
	"time"
)

func train(seriesLength int, dataMap map[model.Ticker][]*model.Period) error {

	log.Infof("Start train of periods with length %v...", seriesLength)
	startTime := time.Now()

	var trainErrors error

	var tickers []model.Ticker
	for ticker, periods := range dataMap {

		if len(periods) < 2 {
			trainErrors =
				multierror.Append(errors.New("unable to train: period sequence must have at least 2 periods"))
			continue
		}

		tickers = append(tickers, ticker)

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

			seqResult := model.Calc(prev.Close, period.Close)
			period.DailyResult = seqResult
			// This period become the previous period
			prev = period
		}

		patterns, err := trainSeries(seriesLength, periods)
		if err != nil {
			trainErrors = multierror.Append(trainErrors, err)
			continue
		}

		err = Repos.PeriodRepo.InsertMany(periods)
		if err != nil {
			trainErrors = multierror.Append(trainErrors, err)
		}

		err = Repos.PatternRepo.InsertMany(patterns)
		if err != nil {
			trainErrors = multierror.Append(trainErrors, err)
		}

		series := &model.Series{Symbol: ticker.Symbol, Length: seriesLength,
			Name: strconv.Itoa(seriesLength) + "-period-series"}
		err = Repos.SeriesRepo.InsertOne(series)
		if err != nil {
			trainErrors = multierror.Append(trainErrors, err)
		}
	}

	err := Repos.TickerRepo.InsertMany(tickers)

	if err != nil {
		trainErrors = multierror.Append(trainErrors, err)
		log.Infof("Completed train of period with length %v with errors took %0.2f minutes",
			seriesLength, time.Since(startTime).Minutes())
		return trainErrors
	}

	log.Infof("Success training periods with length %v with errors took %0.2f minutes",
		seriesLength, time.Since(startTime).Minutes())

	return nil
}

func trainSeries(seriesLength int, periods []*model.Period) ([]*model.Pattern, error) {

	var patterns []*model.Pattern
	if len(periods) < seriesLength+1 {
		return patterns, fmt.Errorf("unable to train series: a series length of %v, needs at least %v periods",
			seriesLength, seriesLength+1)
	}

	var patternMap = make(map[string]*model.Pattern)
	for i, period := range periods {

		// Skip until we have enough testInputData the pattern sequence
		// Must have at least series length + 1 to train
		if i <= seriesLength {
			continue
		}

		// Previous pattern name, such as 'UUD' for a pattern of Up -> Up -> Down.
		var patName string
		for x := seriesLength; x >= 1; x-- {
			patName += fmt.Sprint(periods[i-x].DailyResult)
		}
		r := model.Calc(periods[i-1].Close, period.Close)

		// Find the pattern and increment the total for the given result, r
		var pattern *model.Pattern
		pattern, found := patternMap[patName]
		if !found {
			pattern = &model.Pattern{}
			pattern.Symbol = period.Symbol
			pattern.Value = patName
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

	for _, v := range patternMap {
		patterns = append(patterns, v)
	}

	return patterns, nil
}
