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

type PredictTestSuite struct {
	suite.Suite
}

func TestPredictTestSuite(t *testing.T) {
	suite.Run(t, new(PredictTestSuite))
}

func (suite *PredictTestSuite) SetupTest() {
	conf := config.Init("app-config-test.yaml")
	Repos = mal.New(conf)
}

func (suite *PredictTestSuite) TearDownTest() {
	Repos.DropAll(suite.T())
}

func (suite *PredictTestSuite) TestPredict() {

	r := csv.NewReader(strings.NewReader(testInputData))
	r.TrimLeadingSpace = true
	dataMap := make(map[model.Ticker][]*model.Period)
	err := loadData("test", r, testCompanyData, dataMap)
	assert.NoError(suite.T(), err)

	err = train(3, dataMap)
	assert.NoError(suite.T(), err)

	prediction, err := predict("test")
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), "test", prediction.TickerSymbol)
	assert.Equal(suite.T(), 0.5, prediction.Series[0].ProbabilityUp)
	assert.Equal(suite.T(), 0.5, prediction.Series[0].ProbabilityDown)
}
