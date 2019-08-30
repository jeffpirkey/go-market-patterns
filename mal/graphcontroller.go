package mal

import (
	"go-market-patterns/model/graph"
	"strings"
)

type GraphController struct {
	periodRepo  PeriodRepo
	patternRepo PatternRepo
}

func (c *GraphController) FindPeriodCloseSeries(symbol string) ([]interface{}, error) {

	slice, err := c.periodRepo.FindBySymbol(symbol, SortAsc)
	if err != nil {
		return nil, err
	}

	size := len(slice)
	var ohlc = make([][]interface{}, size)
	var vol = make([][]interface{}, size)
	for i, period := range slice {
		ohlc[i] = make([]interface{}, 5)
		date := period.Date.UnixNano() / 1000000
		ohlc[i][0] = date
		ohlc[i][1] = period.Open
		ohlc[i][2] = period.High
		ohlc[i][3] = period.Low
		ohlc[i][4] = period.Close
		vol[i] = make([]interface{}, 2)
		vol[i][0] = date
		vol[i][1] = period.Volume
	}

	stockPrice := graph.StockSeries{"ohlc", symbol + "-ohlc",
		strings.ToUpper(symbol) + " Stock Price", ohlc, 0}
	//stockVol := graph.StockSeries{"column", symbol + "-volume", strings.ToUpper(symbol) + " Volume", vol, 1}
	series := make([]interface{}, 1)
	series[0] = stockPrice
	//series[1] = stockVol

	return series, nil
}

func (c *GraphController) FindPatternDensities(symbol string) (*graph.PatternDensity, error) {

	patterns, err := c.patternRepo.FindBySymbol(symbol)
	if err != nil {
		return nil, err
	}

	data := graph.PatternDensity{}
	size := len(patterns)
	data.Categories = make([]string, size)
	data.Totals = make([]int, size)
	data.Ups = make([]int, size)
	data.Downs = make([]int, size)
	data.NoChanges = make([]int, size)
	idx := 0
	for _, pattern := range patterns {
		data.Categories[idx] = pattern.Value
		data.Ups[idx] = pattern.UpCount
		data.Downs[idx] = pattern.DownCount
		data.NoChanges[idx] = pattern.NoChangeCount
		data.Totals[idx] = pattern.TotalCount
		idx++
	}

	return &data, nil
}
