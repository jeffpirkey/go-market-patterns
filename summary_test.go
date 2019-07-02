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
	//assert := assert.New(t)

	csvFile, _ := os.Open("data/ibm.us.txt")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	load("ibm", reader)

	train("ibm")

	summary("ibm")
}

func TestFind50File(t *testing.T) {
	assert := assert.New(t)

	csvFile, _ := os.Open("data/ibm.us.txt")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	load("ibm", reader)

	train("ibm")

	found := find50("ibm")

	assert.NotEmpty(found)
}

func TestFindLast(t *testing.T) {
	assert := assert.New(t)

	in := strings.Join([]string{
		`date,value`,
		`2019-01-01, 2`,
		`2019-01-02, 3`,
		`2019-01-03, 4`,
		`2019-01-04, 5`,
		`2019-01-05, 2`,
		`2019-01-06, 3`,
		`2019-01-07, 4`,
		`2019-01-08, 5`,
		`2019-01-09, 6`,
		`2019-01-10, 2`,
		`2019-01-11, 3`,
		`2019-01-13, 4`,
		`2019-01-14, 5`,
		`2019-01-15, 5`,
	}, "\n")

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	load("test", r)

	train("test")

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
