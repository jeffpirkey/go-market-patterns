package main

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"market-patterns/model"
)

func trainAll() error {

	var results error
	for _, symbol := range *Repos.TickerRepo.FindSymbols() {
		err := train(symbol)
		if err != nil {
			results = multierror.Append(results, errors.Wrapf(err, "error training %v", symbol))
		}
	}
	return results
}

func train(symbol string) error {

	ticker := Repos.TickerRepo.FindOne(symbol)

	if ticker == nil {
		return errors.Errorf("unable to train ticker %s due to symbol not found in repo", symbol)
	}

	// Get a slice of descending sort of periods by date
	periods := ticker.PeriodSlice()

	if len(periods) < 2 {
		return fmt.Errorf("unable to train: period sequence must have at least 2 periods")
	}

	// Train the day-to-day results between
	// two consecutive periods across our period slice
	var prev *model.Period
	for i, period := range periods {

		// Set the first index to prev and skip,
		// as we can't compare it to anything
		if i == 0 {
			period.SequenceResult = model.NotDefined
			prev = period
			continue
		}

		if prev == nil {
			return errors.New("previous period not set iterating during train")
		}

		seqResult := model.Calc(prev.Close, period.Close)
		period.SequenceResult = seqResult
		// This period become the previous period
		prev = period
	}

	Repos.TickerRepo.FindOneAndReplace(ticker)

	return nil
}

func trainAllSeries(seriesName, seriesDesc string, seriesLen int) error {

	var results error
	for _, symbol := range *Repos.TickerRepo.FindSymbols() {
		err := trainSeries(symbol, seriesName, seriesDesc, seriesLen)
		if err != nil {
			results = multierror.Append(results, errors.Wrapf(err, "error training %v", symbol))
		}
	}

	return results
}

func trainSeries(symbol, seriesName, seriesDesc string, seriesLen int) error {

	ticker := Repos.TickerRepo.FindOne(symbol)

	// Get a slice of descending sort of periods by date
	periods := ticker.PeriodSlice()

	if len(periods) < seriesLen+1 {
		return fmt.Errorf("unable to train series: a series length of %v, needs at least %v periods",
			seriesLen, seriesLen+1)
	}

	// We have a valid series, so we can add it to the ticker
	series := ticker.AddSeries(seriesName, seriesDesc, seriesLen)

	for i, period := range periods {

		// Skip until we have enough in the pattern sequence
		// Must have at least series length + 1 to train
		if i <= seriesLen {
			continue
		}

		// Previous pattern name, such as 'UUD' for a pattern of Up -> Up -> Down.
		var patName string
		for x := series.SeriesLength; x >= 1; x-- {
			patName += fmt.Sprint(periods[i-x].SequenceResult)
		}
		r := model.Calc(periods[i-1].Close, period.Close)

		// Find the pattern and increment the total for the given result, r
		pattern := ticker.FindPattern(patName)
		pattern.Inc(r)

		// Store the result for the series name being trained in the period
		period.AddSeriesResult(seriesName, r)
	}

	Repos.TickerRepo.FindOneAndReplace(ticker)

	return nil
}
