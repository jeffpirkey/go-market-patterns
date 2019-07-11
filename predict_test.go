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

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	err := load("test", r)
	assert.NoError(suite.T(), err, "Expected no errors loading test data")

	err = train("test")
	if err != nil {
		assert.Fail(suite.T(), "test error", err)
	}
	err = trainSeries("test", "3-periods", "3 period series", 3)
	if err != nil {
		assert.Fail(suite.T(), "test error", err)
	}

	prediction, err := predict("test")
	if err != nil {
		assert.Fail(suite.T(), "Did not expect an error", err)
	}
	assert.Equal(suite.T(), "test", prediction.TickerSymbol,
		"Expected ticker system to be equal to 'test'")

	assert.Equal(suite.T(), 0.5, prediction.Series[0].Probabilities[model.Up],
		"Expected prediction for Up to be 50%")
	assert.Equal(suite.T(), 0.5, prediction.Series[0].Probabilities[model.Down],
		"Expected prediction for Up to be 50%")

	//fmt.Println(*utils.ToJsonString(prediction))
}
