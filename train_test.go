package main

import (
	"bufio"
	"encoding/csv"
	"github.com/stretchr/testify/assert"
	"market-patterns/model"
	"os"
	"strings"
	"testing"
)

var in = strings.Join([]string{
	`Date,Open,High,Low,Close,Volume`, // Header skipped
	`2019-01-01, 1, 2, 3, 2, 100`,     // N/A
	`2019-01-02, 1, 2, 3, 3, 101`,     // N/A
	`2019-01-03, 1, 2, 3, 4, 102`,     // N/A
	`2019-01-04, 1, 2, 3, 5, 103`,     // N/A
	`2019-01-05, 1, 2, 3, 2, 104`,     // UUU -> D
	`2019-01-06, 1, 2, 3, 3, 105`,     // UUD -> U
	`2019-01-07, 1, 2, 3, 4, 106`,     // UDU -> U
	`2019-01-08, 1, 2, 3, 5, 107`,     // DUU -> U
	`2019-01-09, 1, 2, 3, 6, 108`,     // UUU -> U
	`2019-01-10, 1, 2, 3, 2, 109`,     // UUU -> D
	`2019-01-11, 1, 2, 3, 3, 110`,     // UUD -> U
	`2019-01-12, 1, 2, 3, 4, 111`,     // UDU -> U
	`2019-01-13, 1, 2, 3, 5, 112`,     // DUU -> U
	`2019-01-14, 1, 2, 3, 6, 113`,     // UUU -> U
}, "\n")

func TestTrain(t *testing.T) {
	assert := assert.New(t)

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	load("test", r)

	err := train("test")
	if err != nil {
		assert.Fail("test error", err)
	}

	ticker := Tickers.Find("test")
	slice := ticker.PeriodSlice()

	assert.Equal(model.NotDefined, slice[0].SequenceResult, "Expected first sequence to be Not Defined")
	assert.Equal(model.Up, slice[1].SequenceResult, "Expected sequence to be Up")
	assert.Equal(model.Down, slice[4].SequenceResult, "Expected sequence to be Down")
	assert.Equal(model.Up, slice[13].SequenceResult, "Expected last sequence to be Up")
	assert.Equal(model.Up, slice.Last().SequenceResult, "Expected last sequence via Last() to be Up")
}

func TestTrainSeries(t *testing.T) {
	assert := assert.New(t)

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	load("test", r)

	err := train("test")
	if err != nil {
		assert.Fail("test error", err)
	}
	err = trainSeries("test", "3-period-series", "3 period series", 3)
	if err != nil {
		assert.Fail("test error", err)
	}

	ticker := Tickers.Find("test")
	assert.NotEmpty(ticker.FindAllPatterns(), "Expected patterns to be populated")
}

func TestTrainFile(t *testing.T) {
	assert := assert.New(t)

	csvFile, _ := os.Open("data/zf.us.txt")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	load("zf", reader)

	err := train("zf")
	if err != nil {
		assert.Fail("test error", err)
	}
	err = trainSeries("zf", "3-period-series", "3 period series", 3)
	if err != nil {
		assert.Fail("test error", err)
	}

	ticker := Tickers.Find("test")
	assert.NotEmpty(ticker.FindAllPatterns(), "Expected patterns to be populated")
}
