package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/stretchr/testify/assert"
	"market-patterns/model"
	"market-patterns/utils"
	"os"
	"strings"
	"testing"
)

func TestPredict(t *testing.T) {
	assert := assert.New(t)

	r := csv.NewReader(strings.NewReader(in))
	r.TrimLeadingSpace = true

	load("test", r)

	err := train("test")
	if err != nil {
		assert.Fail("test error", err)
	}
	err = trainSeries("test", "3-periods", "3 period series", 3)
	if err != nil {
		assert.Fail("test error", err)
	}

	prediction, err := predict("test")
	if err != nil {
		assert.Fail("Did not expect an error", err)
	}
	assert.Equal("test", prediction.TickerSymbol, "Expected ticker system to be equal to 'test'")

	assert.Equal(0.5, prediction.Series["3-periods"].Probabilities[model.Up.String()],
		"Expected prediction for Up to be 50%")
	assert.Equal(0.5, prediction.Series["3-periods"].Probabilities[model.Down.String()],
		"Expected prediction for Up to be 50%")

	fmt.Println(*utils.ToJsonString(prediction))
}

func TestPredictFile(t *testing.T) {
	assert := assert.New(t)

	csvFile, _ := os.Open("data/ibm.us.txt")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	load("ibm", reader)
	err := train("ibm")
	if err != nil {
		assert.Fail("test error", err)
	}
	err = trainSeries("ibm", "3-period-series", "3 period series", 3)
	if err != nil {
		assert.Fail("test error", err)
	}

	prediction, err := predict("ibm")
	if err != nil {
		assert.Fail("Did not expect an error", err)
	}
	assert.Equal("ibm", prediction.TickerSymbol, "Expected ticker system to be equal to 'ibm'")

	fmt.Println(*utils.ToJsonString(prediction))
}
