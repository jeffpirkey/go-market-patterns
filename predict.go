package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"market-patterns/mal"
	"market-patterns/model"
	"market-patterns/model/report"
	"strings"
	"time"
)

func predict(symbol string) (report.Prediction, error) {

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

		if strings.Contains(match, model.NotDefined) {
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
