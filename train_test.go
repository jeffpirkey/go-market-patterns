package main

import (
	"encoding/csv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"market-patterns/config"
	"market-patterns/mal"
	"market-patterns/model"
	"strings"
	"testing"
)

type TrainTestSuite struct {
	suite.Suite
}

func TestTrainTestSuite(t *testing.T) {
	suite.Run(t, new(TrainTestSuite))
}

func (suite *TrainTestSuite) SetupTest() {
	conf := config.Init("app-config-test.yaml")
	Repos = mal.New(conf)
}

func (suite *TrainTestSuite) TearDownTest() {
	Repos.DropAll(suite.T())
}

// *********************************************************
// 	 Test Train functions
// *********************************************************

func (suite *TrainTestSuite) TestTrain() {

	r := csv.NewReader(strings.NewReader(testInputData))
	r.TrimLeadingSpace = true

	dataMap := make(map[model.Ticker][]*model.Period)
	err := loadData("test", r, testCompanyData, dataMap)
	assert.NoError(suite.T(), err)

	err = train(3, dataMap)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), dataMap)

	var periods model.PeriodSlice
	for _, v := range dataMap {
		periods = v
	}
	assert.Equal(suite.T(), model.NotDefined, periods[0].DailyResult)
	assert.Equal(suite.T(), model.Up, periods[1].DailyResult)
	assert.Equal(suite.T(), model.Down, periods[4].DailyResult)
	assert.Equal(suite.T(), model.Up, periods[13].DailyResult)
	assert.Equal(suite.T(), model.Up, periods.Last().DailyResult)
}

func (suite *TrainTestSuite) TestTrainBadPeriodSize() {

	r := csv.NewReader(strings.NewReader(testInBadPeriodLength))
	r.TrimLeadingSpace = true
	dataMap := make(map[model.Ticker][]*model.Period)
	err := loadData("test", r, testCompanyData, dataMap)
	assert.NoError(suite.T(), err)

	err = train(3, dataMap)
	assert.Error(suite.T(), err)
}
