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

func TestSummaryFile(t *testing.T) {
	assert := assert.New(t)

	csvFile, _ := os.Open("data/ibm.us.txt")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	load("ibm", reader)
	err := train("ibm")
	if err != nil {
		assert.Fail("test error", err)
	}
	err = trainSeries("ibm", "3-period-series", "3 period series", 3)
	if err != nil {
		assert.Fail("test error", err)
	}

	summary("ibm")
}

func TestFind50File(t *testing.T) {
	assert := assert.New(t)

	csvFile, _ := os.Open("data/ibm.us.txt")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	load("ibm", reader)
	err := train("ibm")
	if err != nil {
		assert.Fail("test error", err)
	}
	err = trainSeries("ibm", "3-period-series", "3 period series", 3)
	if err != nil {
		assert.Fail("test error", err)
	}
	found := find50("ibm")

	assert.NotEmpty(found)
}

func TestFindLast(t *testing.T) {
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

	last, err := findLastPeriod("test")

	if err != nil {
		assert.Fail(error.Error(err))
	}

	expected, err := time.Parse(timeFormat, "2019-01-15")

	if err != nil {
		assert.Fail(error.Error(err))
	}

	assert.Equal(expected, last.Date, "Expected last date to be equal")
}
