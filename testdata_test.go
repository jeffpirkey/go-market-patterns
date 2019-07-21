package main

import "strings"

var testInputData = strings.Join([]string{
	`Date,Open,High,Low,Close,Volume`, // Header skipped
	`2019-01-01, 1, 2, 3, 2, 100`,     // N/A
	`2019-01-02, 1, 2, 3, 3, 101`,     // N/A
	`2019-01-03, 1, 2, 3, 4, 102`,     // N/A
	`2019-01-04, 1, 2, 3, 5, 103`,     // N/A
	`2019-01-05, 1, 2, 3, 2, 104`,     // UUU -> D
	`2019-01-06, 1, 2, 3, 3, 105`,     // UUD -> U
	`2019-01-07, 1, 2, 3, 4, 106`,     // UDU -> U
	`2019-01-08, 1, 2, 3, 5, 107`,     // DUU -> U
	`2019-01-09, 1, 2, 3, 6, 108`,     // UUU -> U
	`2019-01-10, 1, 2, 3, 2, 109`,     // UUU -> D
	`2019-01-11, 1, 2, 3, 3, 110`,     // UUD -> U
	`2019-01-12, 1, 2, 3, 4, 111`,     // UDU -> U
	`2019-01-13, 1, 2, 3, 5, 112`,     // DUU -> U
	`2019-01-14, 1, 2, 3, 6, 113`,     // UUU -> U
}, "\n")

var testInBadPeriodLength = strings.Join([]string{
	`Date,Open,High,Low,Close,Volume`, // Header skipped
	`2019-01-01, 1, 2, 3, 2, 100`,     // N/A
}, "\n")

var testCompanyData = map[string]string{"test": "Test Company"}
var testBadCompanyData = map[string]string{"bad": "Bad"}

var testIbmFile = "data/test/ibm.us.txt"
var testCompanyFile = "data/nyse-symb-name.csv"
