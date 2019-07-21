package main

import (
	"bufio"
	"encoding/csv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"market-patterns/config"
	"market-patterns/mal"
	"market-patterns/model"
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
	dataMap := make(map[model.Ticker][]*model.Period)
	err := loadData("test", r, testCompanyData, dataMap)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), dataMap)
}

func (suite *LoadCsvTestSuite) TestLoadBadCompanyFile() {
	dataMap := make(map[model.Ticker][]*model.Period)
	err := load("blah", "badcompany", dataMap)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestLoadBadUrl() {
	dataMap := make(map[model.Ticker][]*model.Period)
	err := load("blah", testCompanyFile, dataMap)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestLoadEmptyDataFile() {
	dataMap := make(map[model.Ticker][]*model.Period)
	err := load("data/test/empty.txt", testCompanyFile, dataMap)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestLoadEmptyCompanyFile() {
	_, err := loadCompanies("data/test/empty.txt")
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestBadDateLoad() {
	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`bad, 2, 3, 1, 2, 100`,
	}, "\n")
	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true
	dataMap := make(map[model.Ticker][]*model.Period)
	err := loadData("test", r, testCompanyData, dataMap)

	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestBadOpenLoad() {
	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-01, bad, 3, 1, 2, 100`,
	}, "\n")
	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	dataMap := make(map[model.Ticker][]*model.Period)
	err := loadData("test", r, testCompanyData, dataMap)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestBadHighLoad() {
	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-01, 2, bad, 1, 2, 100`,
	}, "\n")
	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	dataMap := make(map[model.Ticker][]*model.Period)
	err := loadData("test", r, testCompanyData, dataMap)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestBadLowLoad() {
	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-01, 2, 3, bad, 2, 100`,
	}, "\n")
	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	dataMap := make(map[model.Ticker][]*model.Period)
	err := loadData("test", r, testCompanyData, dataMap)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestBadCloseLoad() {
	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-01, 2, 3, 1, bad, 100`,
	}, "\n")
	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	dataMap := make(map[model.Ticker][]*model.Period)
	err := loadData("test", r, testCompanyData, dataMap)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestVolumeLoad() {
	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-01, 2, 3, 1, 2, bad`,
	}, "\n")
	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	dataMap := make(map[model.Ticker][]*model.Period)
	err := loadData("test", r, testCompanyData, dataMap)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestLoadSortUnordered() {

	date1, err := time.Parse(timeFormat, "2019-01-01")
	assert.NoError(suite.T(), err)

	date2, err := time.Parse(timeFormat, "2019-01-02")
	assert.NoError(suite.T(), err)

	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-02, 2, 3, 1, 2, 100`,
		`2019-01-01, 3, 4, 1, 4, 101`,
	}, "\n")
	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	dataMap := make(map[model.Ticker][]*model.Period)
	err = loadData("test", r, testCompanyData, dataMap)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 1, len(dataMap))

	var ary []*model.Period
	for _, a := range dataMap {
		assert.Equal(suite.T(), 2, len(a))
		ary = a
		break
	}

	first := ary[0]
	second := ary[1]

	assert.Equal(suite.T(), date1, first.Date)
	assert.Equal(suite.T(), date2, second.Date)
}

func (suite *LoadCsvTestSuite) TestLoadSortOrdered() {

	date1, err := time.Parse(timeFormat, "2019-01-01")
	assert.NoError(suite.T(), err)

	date2, err := time.Parse(timeFormat, "2019-01-02")
	assert.NoError(suite.T(), err)

	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-01, 2, 3, 1, 2, 100`,
		`2019-01-02, 3, 4, 1, 4, 101`,
	}, "\n")
	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	dataMap := make(map[model.Ticker][]*model.Period)
	err = loadData("test", r, testCompanyData, dataMap)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), 1, len(dataMap))

	var ary []*model.Period
	for _, a := range dataMap {
		assert.Equal(suite.T(), 2, len(a))
		ary = a
		break
	}

	first := ary[0]
	second := ary[1]
	assert.Equal(suite.T(), date1, first.Date)
	assert.Equal(suite.T(), date2, second.Date)
}

func (suite *LoadCsvTestSuite) TestLoadInvalidFile() {

	csvFile, _ := os.Open("data/test-exceptions/noexists.txt")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	dataMap := make(map[model.Ticker][]*model.Period)
	err := loadData("bad", reader, testBadCompanyData, dataMap)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestLoadInvalidDir() {

	dataMap := make(map[model.Ticker][]*model.Period)
	companyData, err := loadCompanies("data/nyse-symb-name.csv")
	assert.NoError(suite.T(), err)
	err = loadDir("data/invalid/", companyData, dataMap)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestLoadDirWithInvalidCSV() {

	dataMap := make(map[model.Ticker][]*model.Period)
	companyData, err := loadCompanies("data/nyse-symb-name.csv")
	assert.NoError(suite.T(), err)
	err = loadDir("data/test-dir/", companyData, dataMap)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestLoadInvalidZipArchive() {

	dataMap := make(map[model.Ticker][]*model.Period)
	companyData, err := loadCompanies("data/nyse-symb-name.csv")
	assert.NoError(suite.T(), err)
	err = loadZip("data/test/invalid.zip", companyData, dataMap)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestLoadZipArchiveWithInvalidCSV() {

	dataMap := make(map[model.Ticker][]*model.Period)
	companyData, err := loadCompanies("data/nyse-symb-name.csv")
	assert.NoError(suite.T(), err)
	err = loadZip("data/test-dir/empty.txt.zip", companyData, dataMap)
	assert.Error(suite.T(), err)
}
