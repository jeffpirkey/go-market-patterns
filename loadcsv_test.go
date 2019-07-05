package main

import (
	"bufio"
	"encoding/csv"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	assert := assert.New(t)

	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-01, 2, 3, 1, 2, 100`,
		`2019-01-02, 3, 4, 1, 4, 101`,
	}, "\n")

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	load("test", r)

	ticker := Tickers.Find("test")
	periods := ticker.FindAllPeriods()
	assert.NotEmpty(periods, "Periods is not empty after load")
	assert.Equal(2, len(periods), "Expected Periods have 1 entry")
}

func TestLoadSortUnordered(t *testing.T) {
	assert := assert.New(t)

	date1, err := time.Parse(timeFormat, "2019-01-01")
	if err != nil {
		assert.Fail("failed parse date1")
	}

	date2, err := time.Parse(timeFormat, "2019-01-02")
	if err != nil {
		assert.Fail("failed parse date2")
	}

	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-02, 2, 3, 1, 2, 100`,
		`2019-01-01, 3, 4, 1, 4, 101`,
	}, "\n")

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	load("test", r)

	ticker := Tickers.Find("test")
	slice := ticker.PeriodSlice()
	first := slice[0]
	second := slice[1]

	assert.Equal(date1, first.Date, "Expected first data to be the first")
	assert.Equal(date2, second.Date, "Expected second data to be the last")
}

func TestLoadSortOrdered(t *testing.T) {
	assert := assert.New(t)

	date1, err := time.Parse(timeFormat, "2019-01-01")
	if err != nil {
		assert.Fail("failed parse date1")
	}

	date2, err := time.Parse(timeFormat, "2019-01-02")
	if err != nil {
		assert.Fail("failed parse date2")
	}

	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-01, 2, 3, 1, 2, 100`,
		`2019-01-02, 3, 4, 1, 4, 101`,
	}, "\n")

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	load("test", r)

	ticker := Tickers.Find("test")
	slice := ticker.PeriodSlice()
	first := slice[0]
	second := slice[1]

	assert.Equal(date1, first.Date, "Expected first data to be the first")
	assert.Equal(date2, second.Date, "Expected second data to be the last")
}

func TestLoadFile(t *testing.T) {
	assert := assert.New(t)
	csvFile, _ := os.Open("data/zf.us.txt")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	load("zf", reader)
	ticker := Tickers.Find("zf")
	assert.NotEmpty(ticker.FindAllPeriods(), "Expected Periods to be populated")
}

func TestLoadDir(t *testing.T) {
	assert := assert.New(t)

	loadDir("data/")

	assert.NotEmpty(Tickers, "Expected tickers to not be empty")
	assert.Equal(3, len(Tickers.FindAll()), "Expected 3 tickers")
}
