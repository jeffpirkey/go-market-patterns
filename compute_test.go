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

type ComputeTestSuite struct {
	suite.Suite
}

func TestComputeTestSuite(t *testing.T) {
	suite.Run(t, new(ComputeTestSuite))
}

func (suite *ComputeTestSuite) SetupSuite() {
	conf := config.Init()
	Repos = mal.New(conf)
}

func (suite *ComputeTestSuite) TearDownTest() {
	Repos.DropAll(suite.T())
}

// *********************************************************
// 	 Test compute functions
// *********************************************************

func (suite *ComputeTestSuite) TestComputeSeries() {

	r := csv.NewReader(strings.NewReader(testInputData))
	r.TrimLeadingSpace = true

	err := loadAndTrainData("test", "test company", r, 3)
	assert.NoError(suite.T(), err)
}
