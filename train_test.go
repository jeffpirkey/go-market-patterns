package main

import (
	"encoding/csv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go-market-patterns/config"
	"go-market-patterns/mal"
	"strings"
	"testing"
)

type TrainTestSuite struct {
	suite.Suite
}

func TestTrainTestSuite(t *testing.T) {
	suite.Run(t, new(TrainTestSuite))
}

func (suite *TrainTestSuite) SetupSuite() {
	conf := config.Init()
	Repos = mal.New(conf)
}

func (suite *TrainTestSuite) TearDownTest() {
	Repos.DropAll(suite.T())
}

// *********************************************************
// 	 Test Train functions
// *********************************************************

func (suite *TrainTestSuite) TestTrainAllDaily() {

	r := csv.NewReader(strings.NewReader(testInputData))
	r.TrimLeadingSpace = true

	err := loadAndTrainData("test", "test compnay", r, 3)
	assert.NoError(suite.T(), err)

	/*
		assert.Equal(suite.T(), model.NotDefined, periods[0].DailyResult)
		assert.Equal(suite.T(), model.Up, periods[1].DailyResult)
		assert.Equal(suite.T(), model.Down, periods[4].DailyResult)
		assert.Equal(suite.T(), model.Up, periods[13].DailyResult)
		assert.Equal(suite.T(), model.Up, periods.Last().DailyResult)

	*/
}
