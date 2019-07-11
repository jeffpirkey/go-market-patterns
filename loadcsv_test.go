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

type LoadCsvTestSuite struct {
	suite.Suite
}

func TestLoadCsvTestSuite(t *testing.T) {
	suite.Run(t, new(LoadCsvTestSuite))
}

func (suite *LoadCsvTestSuite) SetupTest() {
	conf := config.Init("app-config-test.yaml")
	Repos = mal.New(conf)
}

func (suite *LoadCsvTestSuite) TearDownTest() {
	Repos.DropAll(suite.T())
}

func (suite *LoadCsvTestSuite) TestLoad() {

	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-01, 2, 3, 1, 2, 100`,
		`2019-01-02, 3, 4, 1, 4, 101`,
	}, "\n")

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	err := load("test", r)
	assert.NoError(suite.T(), err, "Expected no errors load test")

	ticker := Repos.TickerRepo.FindOne("test")
	periods := ticker.FindAllPeriods()
	assert.NotEmpty(suite.T(), periods, "Periods is not empty after load")
	assert.Equal(suite.T(), 2, len(periods), "Expected Periods have 2 entries")
}

func (suite *LoadCsvTestSuite) TestBadDateLoad() {
	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`bad, 2, 3, 1, 2, 100`,
	}, "\n")
	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true
	err := load("test", r)
	assert.Error(suite.T(), err, "Expected bad date error")
}

func (suite *LoadCsvTestSuite) TestBadOpenLoad() {
	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-01, bad, 3, 1, 2, 100`,
	}, "\n")
	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true
	err := load("test", r)
	assert.Error(suite.T(), err, "Expected bad open error")
}

func (suite *LoadCsvTestSuite) TestBadHighLoad() {
	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-01, 2, bad, 1, 2, 100`,
	}, "\n")
	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true
	err := load("test", r)
	assert.Error(suite.T(), err, "Expected bad high error")
}

func (suite *LoadCsvTestSuite) TestBadLowLoad() {
	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-01, 2, 3, bad, 2, 100`,
	}, "\n")
	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true
	err := load("test", r)
	assert.Error(suite.T(), err, "Expected bad low error")
}

func (suite *LoadCsvTestSuite) TestBadCloseLoad() {
	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-01, 2, 3, 1, bad, 100`,
	}, "\n")
	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true
	err := load("test", r)
	assert.Error(suite.T(), err, "Expected bad close error")
}

func (suite *LoadCsvTestSuite) TestVolumeLoad() {
	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-01, 2, 3, 1, 2, bad`,
	}, "\n")
	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true
	err := load("test", r)
	assert.Error(suite.T(), err, "Expected bad volume error")
}

func (suite *LoadCsvTestSuite) TestLoadSortUnordered() {

	date1, err := time.Parse(timeFormat, "2019-01-01")
	assert.NoError(suite.T(), err, "Expected time parsing to not have nay errors")

	date2, err := time.Parse(timeFormat, "2019-01-02")
	if err != nil {
		assert.Fail(suite.T(), "failed parse date2")
	}

	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-02, 2, 3, 1, 2, 100`,
		`2019-01-01, 3, 4, 1, 4, 101`,
	}, "\n")

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	err = load("test", r)
	assert.NoError(suite.T(), err, "Expected no errors load test")

	ticker := Repos.TickerRepo.FindOne("test")
	slice := ticker.PeriodSlice()
	first := slice[0]
	second := slice[1]

	assert.Equal(suite.T(), date1, first.Date, "Expected first data to be the first")
	assert.Equal(suite.T(), date2, second.Date, "Expected second data to be the last")
}

func (suite *LoadCsvTestSuite) TestLoadSortOrdered() {

	date1, err := time.Parse(timeFormat, "2019-01-01")
	assert.NoError(suite.T(), err, "Expected time parsing to not have nay errors")

	date2, err := time.Parse(timeFormat, "2019-01-02")
	assert.NoError(suite.T(), err, "Expected time parsing to not have nay errors")

	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-01, 2, 3, 1, 2, 100`,
		`2019-01-02, 3, 4, 1, 4, 101`,
	}, "\n")

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	err = load("test", r)
	assert.NoError(suite.T(), err, "Expected no errors load test")

	ticker := Repos.TickerRepo.FindOne("test")
	slice := ticker.PeriodSlice()
	first := slice[0]
	second := slice[1]

	assert.Equal(suite.T(), date1, first.Date, "Expected first data to be the first")
	assert.Equal(suite.T(), date2, second.Date, "Expected second data to be the last")
}

func (suite *LoadCsvTestSuite) TestLoadFile() {
	csvFile, _ := os.Open("data/test/ibm.us.txt")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	err := load("ibm", reader)
	assert.NoError(suite.T(), err, "Expected no errors load test")

	ticker := Repos.TickerRepo.FindOne("ibm")

	assert.Equal(suite.T(), "ibm", ticker.Symbol, "Expected Ticker symbol to be 'ibm'")
	assert.Equal(suite.T(), 14059, len(ticker.Periods), "Expected Period count to be equal")
	assert.Equal(suite.T(), 0, len(ticker.Patterns), "Expected Period count to be equal")
	assert.Equal(suite.T(), 0, len(ticker.Series), "Expected Period count to be equal")
}

func (suite *LoadCsvTestSuite) TestLoadInvalidFile() {
	csvFile, _ := os.Open("data/test-exceptions/noexists.txt")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	err := load("bad", reader)
	assert.Error(suite.T(), err, "Expected an error loading an invalid file path")
}

func (suite *LoadCsvTestSuite) TestLoadDir() {

	err := loadDir("data/test/")
	assert.NoError(suite.T(), err, "Expected directory load to not have nay errors")

	ticker := Repos.TickerRepo.FindOne("ibm")
	assert.Equal(suite.T(), "ibm", ticker.Symbol, "Expected Ticker symbol to be 'ibm'")
	assert.Equal(suite.T(), 14059, len(ticker.Periods), "Expected Period count to be equal")
	assert.Equal(suite.T(), 0, len(ticker.Patterns), "Expected Period count to be equal")
	assert.Equal(suite.T(), 0, len(ticker.Series), "Expected Period count to be equal")
}

func (suite *LoadCsvTestSuite) TestLoadInvalidDir() {
	err := loadDir("data/invalid/")
	assert.Error(suite.T(), err, "Expected an error loading an invalid directory path")
}

func (suite *LoadCsvTestSuite) TestLoadDirWithInvalidCSV() {
	err := loadDir("data/test-dir/")
	assert.Error(suite.T(), err, "Expected an error loading an invalid directory path")
}

func (suite *LoadCsvTestSuite) TestLoadZipArchive() {

	err := loadZip("data/test/stocks-small.zip")

	assert.NoError(suite.T(), err, "Expected zip archive load to not have errors")

	ticker := Repos.TickerRepo.FindOne("ibm")

	assert.Equal(suite.T(), "ibm", ticker.Symbol, "Expected Ticker symbol to be 'ibm'")
	assert.Equal(suite.T(), 14059, len(ticker.Periods), "Expected Period count to be equal")
	assert.Equal(suite.T(), 0, len(ticker.Patterns), "Expected Period count to be equal")
	assert.Equal(suite.T(), 0, len(ticker.Series), "Expected Period count to be equal")
}

func (suite *LoadCsvTestSuite) TestLoadInvalidZipArchive() {
	err := loadZip("data/test/invalid.zip")
	assert.Error(suite.T(), err, "Expected zip archive load to have errors")
}

func (suite *LoadCsvTestSuite) TestLoadZipArchiveWithInvalidCSV() {
	err := loadZip("data/test-dir/empty.txt.zip")
	assert.Error(suite.T(), err, "Expected zip archive load to have errors")
}
