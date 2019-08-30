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

type PredictTestSuite struct {
	suite.Suite
}

func TestPredictTestSuite(t *testing.T) {
	suite.Run(t, new(PredictTestSuite))
}

func (suite *PredictTestSuite) SetupSuite() {
	conf := config.Init("runtime-config-test.yaml")
	Repos = mal.New(conf)
}

func (suite *PredictTestSuite) TearDownTest() {
	Repos.DropAll(suite.T())
}

func (suite *PredictTestSuite) TestPredict() {

	r := csv.NewReader(strings.NewReader(testInputData))
	r.TrimLeadingSpace = true
	err := loadAndTrainData("test", "test company", r, 3)
	assert.NoError(suite.T(), err)
	prediction, err := predict("test")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test", prediction.TickerSymbol)
	assert.Equal(suite.T(), 0.5, prediction.Series[0].ProbabilityUp)
	assert.Equal(suite.T(), 0.5, prediction.Series[0].ProbabilityDown)
}
