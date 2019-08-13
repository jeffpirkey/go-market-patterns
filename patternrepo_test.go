package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"market-patterns/config"
	"market-patterns/mal"
	"market-patterns/model"
	"testing"
)

type PatternRepoTestSuite struct {
	suite.Suite
}

func TestPatternRepoTestSuite(t *testing.T) {
	suite.Run(t, new(PatternRepoTestSuite))
}

func (suite *PatternRepoTestSuite) SetupSuite() {
	conf := config.Init("runtime-config-test.yaml")
	Repos = mal.New(conf)
}

func (suite *PatternRepoTestSuite) TearDownTest() {
	Repos.DropAll(suite.T())
}

func (suite *PatternRepoTestSuite) TestHighestUp() {

	err := truncAndLoad("data/test/stocks-test.zip", testCompanyFile, 3)
	assert.NoError(suite.T(), err)
	pattern, err := Repos.PatternRepo.FindHighestUpProbability(model.PatternDensityLow)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), pattern)
}
