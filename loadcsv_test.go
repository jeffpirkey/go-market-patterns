package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	assert := assert.New(t)

	in := strings.Join([]string{
		`date,value`,
		`2019-01-01, 2`,
		`2019-01-02, 3`,
	}, "\n")

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	load("test", r)

	assert.NotEmpty(Periods, "Periods is not empty after load")
	assert.Equal(1, len(Periods), "Expected Periods have 1 entry")
	assert.Equal(2, len(Periods["test"]), "Expected Periods for 'test' have 2 entry")
}

func TestLoadSortUnordered(t *testing.T) {
	assert := assert.New(t)

	date1, err := time.Parse(timeFormat, "2019-01-01")
	if err != nil {
		assert.Fail("failed parse date1")
	}

	date2, err := time.Parse(timeFormat, "2019-01-02")
	if err != nil {
		assert.Fail("failed parse date2")
	}

	in := strings.Join([]string{
		`date,value`,
		`2019-01-02, 2`,
		`2019-01-01, 3`,
	}, "\n")

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	load("test", r)

	first := Periods["test"][0]
	second := Periods["test"][1]

	assert.Equal(date1, first.Date, "Expected first data to be the first")
	assert.Equal(date2, second.Date, "Expected second data to be the last")
}

func TestLoadSortOrdered(t *testing.T) {
	assert := assert.New(t)

	date1, err := time.Parse(timeFormat, "2019-01-01")
	if err != nil {
		assert.Fail("failed parse date1")
	}

	date2, err := time.Parse(timeFormat, "2019-01-02")
	if err != nil {
		assert.Fail("failed parse date2")
	}

	in := strings.Join([]string{
		`date,value`,
		`2019-01-01, 2`,
		`2019-01-02, 3`,
	}, "\n")

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	load("test", r)

	first := Periods["test"][0]
	second := Periods["test"][1]

	assert.Equal(date1, first.Date, "Expected first data to be the first")
	assert.Equal(date2, second.Date, "Expected second data to be the last")
}

func TestLoadFile(t *testing.T) {
	//assert := assert.New(t)

	csvFile, _ := os.Open("data/zf.us.txt")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	load("zf", reader)

	fmt.Println(len(Periods["zf"]))

}
