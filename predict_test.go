package main

import (
	"encoding/csv"
	"fmt"
	"strings"
	"testing"
)

func TestPredict(t *testing.T) {
	//assert := assert.New(t)

	in := strings.Join([]string{
		`date,value`,
		`2019-01-01, 2`,
		`2019-01-02, 3`,
		`2019-01-03, 4`,
		`2019-01-04, 5`,
		`2019-01-05, 2`,
		`2019-01-06, 3`,
		`2019-01-07, 4`,
		`2019-01-08, 5`,
		`2019-01-09, 6`,
		`2019-01-10, 2`,
		`2019-01-11, 3`,
		`2019-01-13, 4`,
		`2019-01-14, 5`,
		`2019-01-15, 5`,
	}, "\n")

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	load("test", r)

	train("test")

	err := predict("test", "UUU")

	if err != nil {
		fmt.Println(err)
	}
}
