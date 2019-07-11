package main

import (
	"archive/zip"
	"bufio"
	"encoding/csv"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"market-patterns/model"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	timeFormat = "2006-01-02"
)

func loadDir(path string) error {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return errors.Wrap(err, "unable to load dir")
	}

	var results error

	for _, file := range files {
		split := strings.Split(file.Name(), ".")

		ext := split[len(split)-1]
		if ext != "txt" && ext != "csv" {
			log.Warnf("Skipping unrecognized file extension of %v", ext)
			continue
		}
		// Skip this error and let the load return if the reader is invalid
		csvFile, _ := os.Open(path + file.Name())
		reader := csv.NewReader(bufio.NewReader(csvFile))
		err := load(split[0], reader)
		if err != nil {
			results = multierror.Append(results, err)
		}
	}

	return results
}

func loadZip(zipFile string) error {

	// Open a zip archive for reading.
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return errors.Wrap(err, "problem open zip archive")
	}

	var results error

	defer func(r *zip.ReadCloser) {
		err := r.Close()
		if err != nil {
			results = multierror.Append(results, errors.Wrap(err, "problem closing zip reader"))
		}
	}(r)

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		names := strings.Split(f.Name, ".")
		rc, err := f.Open()
		if err != nil {
			return errors.Wrap(err, "problem open zip file")
		}

		reader := csv.NewReader(rc)
		err = load(names[0], reader)
		if err != nil {
			results = multierror.Append(results, err)
		}
		err = rc.Close()
		if err != nil {
			results = multierror.Append(results, errors.Wrap(err, "problem closing zip file reader"))
		}
	}

	return results
}

func load(symbol string, r *csv.Reader) error {

	var results error

	vals, err := r.ReadAll()
	if err != nil {
		return errors.Wrap(err, "error reading csv")
	}

	if vals == nil {
		return errors.New("Empty or invalid CSV")
	}

	ticker := model.Ticker{Symbol: symbol, Series: make(map[string]*model.Series),
		Patterns: make(map[string]*model.Pattern),
		Periods:  make(map[string]*model.Period)}
	for i, v := range vals {

		if i == 0 {
			// skip header line
			continue
		}

		date, err := convertTime(v[0])
		if err != nil {
			results = multierror.Append(results, errors.Wrap(err, "date field"))
		}
		open, err := convertFloat(v[1])
		if err != nil {
			results = multierror.Append(results, errors.Wrap(err, "open field"))
		}
		high, err := convertFloat(v[2])
		if err != nil {
			results = multierror.Append(results, errors.Wrap(err, "high field"))
		}
		low, err := convertFloat(v[3])
		if err != nil {
			results = multierror.Append(results, errors.Wrap(err, "low field"))
		}
		cl, err := convertFloat(v[4])
		if err != nil {
			results = multierror.Append(results, errors.Wrap(err, "close field"))
		}
		volume, err := convertInt(v[5])
		if err != nil {
			results = multierror.Append(results, errors.Wrap(err, "volume field"))
		}

		p := model.Period{Date: date, Open: open, High: high, Low: low, Close: cl, Volume: volume}
		ticker.AddPeriod(&p)
	}

	// Save to mongo
	Repos.TickerRepo.FindOneAndReplace(&ticker)

	return results
}

func convertFloat(v string) (float64, error) {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return f, errors.Wrap(err, "unable to convert csv value to float")
	}
	return f, nil
}

func convertInt(v string) (int, error) {
	i, err := strconv.Atoi(v)
	if err != nil {
		return i, errors.Wrap(err, "unable to convert csv value to int")
	}
	return i, nil
}

func convertTime(v string) (time.Time, error) {
	t, err := time.Parse(timeFormat, v)
	if err != nil {
		return t, errors.Wrap(err, "unable to convert csv value to time")
	}
	return t, nil
}
