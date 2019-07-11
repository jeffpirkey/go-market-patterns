package main

import (
	"bufio"
	"encoding/csv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"market-patterns/config"
	"market-patterns/mal"
	"os"
	"strings"
	"testing"
	"time"
)

type SummaryTestSuite struct {
	suite.Suite
}

func TestSummaryTestSuite(t *testing.T) {
	suite.Run(t, new(SummaryTestSuite))
}

func (suite *SummaryTestSuite) SetupTest() {
	conf := config.Init("app-config-test.yaml")
	Repos = mal.New(conf)
}

func (suite *LoadCsvTestSuite) SummaryTestSuite() {
	Repos.DropAll(suite.T())
}

func (suite *SummaryTestSuite) TestSummaryFile() {

	csvFile, _ := os.Open("data/test/ibm.us.txt")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	err := load("ibm", reader)
	assert.NoError(suite.T(), err, "Expected no errors loading test data")

	err = train("ibm")
	if err != nil {
		assert.Fail(suite.T(), "test error", err)
	}
	err = trainSeries("ibm", "3-period-series", "3 period series", 3)
	if err != nil {
		assert.Fail(suite.T(), "test error", err)
	}

	summary("ibm")
}

func (suite *SummaryTestSuite) TestFind50File() {

	csvFile, _ := os.Open("data/test/ibm.us.txt")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	err := load("ibm", reader)
	assert.NoError(suite.T(), err, "Expected no errors loading test data")

	err = train("ibm")
	if err != nil {
		assert.Fail(suite.T(), "test error", err)
	}
	err = trainSeries("ibm", "3-period-series", "3 period series", 3)
	if err != nil {
		assert.Fail(suite.T(), "test error", err)
	}
	found := find50("ibm")

	assert.NotEmpty(suite.T(), found)
}

func (suite *SummaryTestSuite) TestFindLast() {

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	err := load("test", r)
	assert.NoError(suite.T(), err, "Expected no errors loading test data")

	err = train("test")
	if err != nil {
		assert.Fail(suite.T(), "test error", err)
	}
	err = trainSeries("test", "3-period-series", "3 period series", 3)
	if err != nil {
		assert.Fail(suite.T(), "test error", err)
	}

	last, err := findLastPeriod("test")

	if err != nil {
		assert.Fail(suite.T(), error.Error(err))
	}

	expected, err := time.Parse(timeFormat, "2019-01-14")

	if err != nil {
		assert.Fail(suite.T(), error.Error(err))
	}

	assert.Equal(suite.T(), expected, last.Date, "Expected last date to be equal")
}
