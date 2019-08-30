package main

import (
	"bufio"
	"encoding/csv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go-market-patterns/config"
	"go-market-patterns/mal"
	"os"
	"strings"
	"testing"
)

type LoadCsvTestSuite struct {
	suite.Suite
}

func TestLoadCsvTestSuite(t *testing.T) {
	suite.Run(t, new(LoadCsvTestSuite))
}

func (suite *LoadCsvTestSuite) SetupSuite() {
	conf := config.Init()
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
	err := loadAndTrainData("test", "test company", r, 0)

	assert.NoError(suite.T(), err)

}

func (suite *LoadCsvTestSuite) TestLoadBadCompanyFile() {
	err := load("blah", "badcompany", 3)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestLoadBadUrl() {
	err := load("blah", testCompanyFile, 3)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestLoadEmptyDataFile() {
	err := load("data/test-exceptions/empty.txt", testCompanyFile, 3)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestLoadEmptyCompanyFile() {
	_, err := loadCompanies("data/test-exceptions/empty.txt")
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestBadDateLoad() {
	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`bad, 2, 3, 1, 2, 100`,
	}, "\n")
	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true
	err := loadAndTrainData("test", "test company", r, 3)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestBadOpenLoad() {
	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-01, bad, 3, 1, 2, 100`,
	}, "\n")
	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true
	err := loadAndTrainData("test", "test company", r, 3)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestBadHighLoad() {
	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-01, 2, bad, 1, 2, 100`,
	}, "\n")
	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true
	err := loadAndTrainData("test", "test company", r, 3)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestBadLowLoad() {
	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-01, 2, 3, bad, 2, 100`,
	}, "\n")
	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	err := loadAndTrainData("test", "test company", r, 3)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestBadCloseLoad() {
	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-01, 2, 3, 1, bad, 100`,
	}, "\n")
	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	err := loadAndTrainData("test", "test company", r, 3)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestVolumeLoad() {
	in := strings.Join([]string{
		`Date,Open,High,Low,Close,Volume`,
		`2019-01-01, 2, 3, 1, 2, bad`,
	}, "\n")
	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	err := loadAndTrainData("test", "test company", r, 3)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestLoadInvalidFile() {

	// Creating a bad reader, so skip the error
	csvFile, _ := os.Open("data/test-exceptions/noexists.txt")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	err := loadAndTrainData("bad", "bad company", reader, 3)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestLoadInvalidDir() {

	companyData, err := loadCompanies("data/nyse-symb-name.csv")
	assert.NoError(suite.T(), err)
	err = loadDir("data/invalid/", companyData, 3)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestLoadDirWithInvalidCSV() {

	companyData, err := loadCompanies("data/nyse-symb-name.csv")
	assert.NoError(suite.T(), err)
	err = loadDir("data/test-dir/", companyData, 3)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestLoadInvalidZipArchive() {

	companyData, err := loadCompanies("data/nyse-symb-name.csv")
	assert.NoError(suite.T(), err)
	err = loadZip("data/test/invalid.zip", companyData, 3)
	assert.Error(suite.T(), err)
}

func (suite *LoadCsvTestSuite) TestLoadZipArchiveWithInvalidCSV() {

	companyData, err := loadCompanies("data/nyse-symb-name.csv")
	assert.NoError(suite.T(), err)
	err = loadZip("data/test-dir/empty.txt.zip", companyData, 3)
	assert.Error(suite.T(), err)
}
