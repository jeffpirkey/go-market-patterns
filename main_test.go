package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"market-patterns/config"
	"market-patterns/mal"
	"testing"
)

type MainTestSuite struct {
	suite.Suite
}

func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}

func (suite *MainTestSuite) SetupTest() {
	conf := config.Init("app-config-test.yaml")
	Repos = mal.New(conf)
}

func (suite *MainTestSuite) TearDownTest() {
	Repos.DropAll(suite.T())
}

func (suite *MainTestSuite) TestTruncLoadFile() {

	err := truncAndLoad(testIbmFile, testCompanyFile)
	assert.NoError(suite.T(), err)

	symbol := "IBM"
	ticker, err := Repos.TickerRepo.FindOne(symbol)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), symbol, ticker.Symbol)

	periods, err := Repos.PeriodRepo.FindBySymbol(symbol, nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 14059, len(periods))

	patterns, err := Repos.PatternRepo.FindBySymbol(symbol)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 26, len(patterns))

	series, err := Repos.SeriesRepo.FindBySymbol(symbol)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(series))
}

func (suite *LoadCsvTestSuite) TestTruncLoadDir() {

	err := truncAndLoad("data/test/", testCompanyFile)
	assert.NoError(suite.T(), err)

	symbol := "IBM"
	ticker, err := Repos.TickerRepo.FindOne(symbol)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), symbol, ticker.Symbol)

	periods, err := Repos.PeriodRepo.FindBySymbol(symbol, nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 14059, len(periods))

	patterns, err := Repos.PatternRepo.FindBySymbol(symbol)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 26, len(patterns))

	series, err := Repos.SeriesRepo.FindBySymbol(symbol)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(series))
}

func (suite *LoadCsvTestSuite) TestTruncLoadZipArchive() {

	err := truncAndLoad("data/test/stocks-test.zip", testCompanyFile)
	assert.NoError(suite.T(), err)

	symbol := "IBM"
	ticker, err := Repos.TickerRepo.FindOne(symbol)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), symbol, ticker.Symbol)

	periods, err := Repos.PeriodRepo.FindBySymbol(symbol, nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 14059, len(periods))

	patterns, err := Repos.PatternRepo.FindBySymbol(symbol)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 26, len(patterns))

	series, err := Repos.SeriesRepo.FindBySymbol(symbol)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, len(series))
}
