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

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	err := load("test", r)
	assert.NoError(suite.T(), err, "Expected no errors loading test data")

	err = train("test")
	assert.NoError(suite.T(), err, "Expect no errors training data")

	ticker := Repos.TickerRepo.FindOne("test")
	slice := ticker.PeriodSlice()

	assert.Equal(suite.T(), model.NotDefined, slice[0].SequenceResult, "Expected first sequence to be Not Defined")
	assert.Equal(suite.T(), model.Up, slice[1].SequenceResult, "Expected sequence to be Up")
	assert.Equal(suite.T(), model.Down, slice[4].SequenceResult, "Expected sequence to be Down")
	assert.Equal(suite.T(), model.Up, slice[13].SequenceResult, "Expected last sequence to be Up")
	assert.Equal(suite.T(), model.Up, slice.Last().SequenceResult, "Expected last sequence via Last() to be Up")
}

func (suite *TrainTestSuite) TestTrainMissingSymbol() {

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true
	err := load("test", r)
	assert.NoError(suite.T(), err, "Expected no errors loading test data")
	err = train("bad")
	assert.Error(suite.T(), err, "Expect train of missing symbol to have errors")
}

func (suite *TrainTestSuite) TestTrainBadPeriodSize() {

	r := csv.NewReader(strings.NewReader(inBadPeriodLength))
	r.TrimLeadingSpace = true
	err := load("test", r)
	assert.NoError(suite.T(), err, "Expected no errors loading test data")
	err = train("test")
	assert.Error(suite.T(), err, "Expect train of missing symbol to have errors")
}

func (suite *TrainTestSuite) TestTrainNoLoad() {
	err := train("test")
	assert.Error(suite.T(), err, "Expect train of missing symbol to have errors")
}

func (suite *TrainTestSuite) TestTrainAll() {

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	err := load("test", r)
	assert.NoError(suite.T(), err, "Expected no errors loading test data")

	err = trainAll()
	assert.NoError(suite.T(), err, "Expect no errors training data")

	ticker := Repos.TickerRepo.FindOne("test")
	slice := ticker.PeriodSlice()

	assert.Equal(suite.T(), model.NotDefined, slice[0].SequenceResult, "Expected first sequence to be Not Defined")
	assert.Equal(suite.T(), model.Up, slice[1].SequenceResult, "Expected sequence to be Up")
	assert.Equal(suite.T(), model.Down, slice[4].SequenceResult, "Expected sequence to be Down")
	assert.Equal(suite.T(), model.Up, slice[13].SequenceResult, "Expected last sequence to be Up")
	assert.Equal(suite.T(), model.Up, slice.Last().SequenceResult, "Expected last sequence via Last() to be Up")
}

func (suite *TrainTestSuite) TestTrainAllWithError() {

	r := csv.NewReader(strings.NewReader(inBadPeriodLength))
	r.TrimLeadingSpace = true

	err := load("test", r)
	assert.NoError(suite.T(), err, "Expected no errors loading test data")

	err = trainAll()
	assert.Error(suite.T(), err, "Expect errors training data")
}

// *********************************************************
// Test train series functions
// *********************************************************

func (suite *TrainTestSuite) TestTrainSeries() {

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	err := load("test", r)
	assert.NoError(suite.T(), err, "Expected no errors loading test data")

	err = train("test")
	assert.NoError(suite.T(), err, "Expect no errors training")

	err = trainSeries("test", "3-period-series", "3 period series", 3)
	assert.NoError(suite.T(), err, "Expect no errors training series")

	ticker := Repos.TickerRepo.FindOne("test")
	assert.NotEmpty(suite.T(), ticker.FindAllPatterns(), "Expected patterns to be populated")
}

func (suite *TrainTestSuite) TestTrainSeriesBadPeriodLength() {

	r := csv.NewReader(strings.NewReader(inBadSeriesLength))
	r.TrimLeadingSpace = true

	err := load("test", r)
	assert.NoError(suite.T(), err, "Expected no errors loading test data")

	err = train("test")
	assert.NoError(suite.T(), err, "Expect no errors training")

	err = trainSeries("test", "3-period-series", "3 period series", 3)
	assert.Error(suite.T(), err, "Expect errors training series")
}

func (suite *TrainTestSuite) TestTrainSeriesAll() {

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	err := load("test", r)
	assert.NoError(suite.T(), err, "Expected no errors loading test data")

	err = trainAll()
	assert.NoError(suite.T(), err, "Expect no errors training data")

	err = trainAllSeries("3-period-series", "3 period series", 3)
	assert.NoError(suite.T(), err, "Expect no errors training series")

	ticker := Repos.TickerRepo.FindOne("test")
	assert.NotEmpty(suite.T(), ticker.FindAllPatterns(), "Expected patterns to be populated")
}

func (suite *TrainTestSuite) TestTrainSeriesAllWithError() {

	r := csv.NewReader(strings.NewReader(inBadSeriesLength))
	r.TrimLeadingSpace = true

	err := load("test", r)
	assert.NoError(suite.T(), err, "Expected no errors loading test data")

	err = trainAll()
	assert.NoError(suite.T(), err, "Expect no errors training")

	err = trainAllSeries("3-period-series", "3 period series", 3)
	assert.Error(suite.T(), err, "Expect error training series")
}
