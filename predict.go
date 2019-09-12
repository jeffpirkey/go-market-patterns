package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-market-patterns/mal"
	"go-market-patterns/model/core"
	"go-market-patterns/model/report"
	"strings"
	"time"
)

func predictOne(symbol string, length int) (report.Prediction, error) {

	log.Infof("Processing prediction for %v with series length %v...", symbol, length)

	var prediction report.Prediction

	series, err := Repos.SeriesRepo.FindOneBySymbolAndLength(symbol, length)
	if err != nil {
		return prediction, err
	}

	slice, err := Repos.PeriodRepo.FindBySymbol(symbol, mal.SortAsc)
	if err != nil {
		return prediction, err
	}

	fromDay := slice.Last().Date
	nextDay := fromDay.AddDate(0, 0, 1)
	prediction = report.Prediction{TickerSymbol: symbol,
		NextDate: fmt.Sprintf("%d-%02d-%02d", nextDay.Year(), nextDay.Month(), nextDay.Day()),
		FromDate: fmt.Sprintf("%d-%02d-%02d", fromDay.Year(), fromDay.Month(), fromDay.Day())}

	lastPeriods := slice.LastByRange(length)
	var match string
	for _, period := range lastPeriods {
		// Find the result for the series name being
		// predicted for each period
		match += period.DailyResult
	}

	ps := report.PredictionSeries{Name: series.Name, Pattern: match}
	prediction.Series = append(prediction.Series, &ps)

	if strings.Contains(match, core.NotDefined) {
		log.Info("No supporting data")
	} else {
		pattern, err := Repos.PatternRepo.FindOneBySymbolAndValueAndLength(symbol, match, length)
		if err != nil {
			return prediction, err
		}

		ps.ProbabilityUp = float64(pattern.UpCount) / float64(pattern.TotalCount)
		ps.ProbabilityDown = float64(pattern.DownCount) / float64(pattern.TotalCount)
		ps.ProbabilityNoChange = float64(pattern.NoChangeCount) / float64(pattern.TotalCount)
	}

	log.Infof("Finished processing prediction for ticker %v and series length %v", symbol, length)

	return prediction, nil
}

func predictAll(symbol string) (report.Prediction, error) {

	startTime := time.Now()

	var prediction report.Prediction
	slice, err := Repos.PeriodRepo.FindBySymbol(symbol, mal.SortAsc)
	if err != nil {
		return prediction, err
	}

	fromDay := slice.Last().Date
	nextDay := fromDay.AddDate(0, 0, 1)
	prediction = report.Prediction{TickerSymbol: symbol,
		NextDate: fmt.Sprintf("%d-%02d-%02d", nextDay.Year(), nextDay.Month(), nextDay.Day()),
		FromDate: fmt.Sprintf("%d-%02d-%02d", fromDay.Year(), fromDay.Month(), fromDay.Day())}

	series, err := Repos.SeriesRepo.FindBySymbol(symbol)
	if err != nil {
		return prediction, err
	}

	for _, s := range series {

		log.Infof("Processing prediction for %v with series %v...", symbol, s.Name)

		lastPeriods := slice.LastByRange(s.Length)
		var match string
		for _, period := range lastPeriods {
			// Find the result for the series name being
			// predicted for each period
			match += period.DailyResult
		}

		ps := report.PredictionSeries{Name: s.Name, Pattern: match}
		prediction.Series = append(prediction.Series, &ps)

		if strings.Contains(match, core.NotDefined) {
			log.Info("No supporting data")
		} else {
			pattern, err := Repos.PatternRepo.FindOneBySymbolAndValueAndLength(symbol, match, s.Length)
			if err != nil {
				return prediction, err
			}

			ps.ProbabilityUp = float64(pattern.UpCount) / float64(pattern.TotalCount)
			ps.ProbabilityDown = float64(pattern.DownCount) / float64(pattern.TotalCount)
			ps.ProbabilityNoChange = float64(pattern.NoChangeCount) / float64(pattern.TotalCount)
		}

		log.Infof("Finished processing prediction for ticker %v and series %v", symbol, s.Name)
	}

	log.Infof("Generating predictions took %0.2f minutes", time.Since(startTime).Minutes())

	return prediction, nil
}
