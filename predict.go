package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"market-patterns/model"
	"market-patterns/report"
	"strings"
)

func predict(tsym string) (report.Prediction, error) {

	ticker := Tickers.Find(tsym)
	slice := ticker.PeriodSlice()

	fromDay := slice.Last().Date
	nextDay := fromDay.AddDate(0, 0, 1)
	prediction := report.Prediction{TickerSymbol: tsym,
		NextDate: fmt.Sprintf("%d-%02d-%02d", nextDay.Year(), nextDay.Month(), nextDay.Day()),
		FromDate: fmt.Sprintf("%d-%02d-%02d", fromDay.Year(), fromDay.Month(), fromDay.Day()),
		Series:   make(map[string]report.PredictionSeries)}

	for seriesName, series := range ticker.FindAllSeries() {

		log.Infof("Processing prediction for %v the series of %v...", tsym, seriesName)

		lastPeriods := slice.LastByRange(series.SeriesLength)
		var match string
		for _, period := range lastPeriods {
			// Find the result for the series name being
			// predicted for each period
			result := period.FindSeriesResult(seriesName)
			match += result.String()
		}

		ps := report.PredictionSeries{Pattern: match, Probabilities: make(map[string]float64)}
		prediction.Series[seriesName] = ps

		if strings.Contains(match, model.NotDefined.String()) {
			log.Info("No supporting data")
		} else {
			pattern := ticker.FindPattern(match)
			for result, count := range pattern.FindAll() {
				pb := float64(count) / float64(pattern.TotalCount())
				ps.Probabilities[result.String()] = pb
			}
		}

		log.Infof("Finished processing prediction for ticker %v and series %v", tsym, seriesName)
	}

	return prediction, nil
}
