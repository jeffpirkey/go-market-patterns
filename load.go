package main

import (
	"archive/zip"
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"market-patterns/model"
	"market-patterns/utils"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	timeFormat = "2006-01-02"
)

func load(url, companyFile string, dataMap map[model.Ticker][]*model.Period) error {

	companyData, err := loadCompanies(companyFile)
	if err != nil {
		return err
	}

	fi, err := os.Stat(url)
	if err != nil {
		return err
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		err = loadDir(url, companyData, dataMap)
	case mode.IsRegular():
		if utils.IsZip(url) {
			err = loadZip(url, companyData, dataMap)
		} else {
			err = loadFile(url, companyData, dataMap)
		}
	}

	return err
}

func loadCompanies(fileName string) (map[string]string, error) {

	log.Infof("Starting load of company data from %v...", fileName)

	startTime := time.Now()

	var data map[string]string

	csvFile, err := os.Open(fileName)
	if err != nil {
		return data, fmt.Errorf("problem loading company data due to %v", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Errorf("unable to close company file due to %v", err)
		}
	}(csvFile)

	reader := csv.NewReader(bufio.NewReader(csvFile))
	vals, err := reader.ReadAll()
	if err != nil {
		return data, errors.Wrapf(err, "error reading %v", fileName)
	}
	if vals == nil {
		return data, fmt.Errorf("empty or invalid CSV in %v", fileName)
	}
	data = make(map[string]string)
	for i, v := range vals {

		if i == 0 {
			// skip header line
			continue
		}
		data[v[0]] = v[1]
	}

	log.Infof("Successful company data load from %v took %0.6f seconds",
		fileName, time.Since(startTime).Seconds())

	return data, nil
}

func loadDir(dataUrl string, companyData map[string]string, dataMap map[model.Ticker][]*model.Period) error {

	log.Infof("Starting load of files from directory %v...", dataUrl)
	startTime := time.Now()

	files, err := ioutil.ReadDir(dataUrl)
	if err != nil {
		return errors.Wrapf(err, "unable to load directory %v", dataUrl)
	}

	var results error

	for _, file := range files {
		split := strings.Split(file.Name(), ".")

		ext := split[len(split)-1]
		if ext != "txt" && ext != "csv" {
			log.Warnf("Skipping unrecognized file extension '%v'", ext)
			continue
		}

		// Skip this error and let the load return if the reader is invalid
		csvFile, _ := os.Open(dataUrl + file.Name())
		reader := csv.NewReader(bufio.NewReader(csvFile))
		// split[0] should be the ticker symbol
		err := loadData(strings.ToUpper(split[0]), reader, companyData, dataMap)
		if err != nil {
			results = multierror.Append(results, err)
		}
		err = csvFile.Close()
		if err != nil {
			results = multierror.Append(results, errors.Wrap(err, "unable to close company file due to %v"))
		}
	}

	if results != nil {
		log.Infof("Completed directory load from %v with errors took %0.2f minutes",
			dataUrl, time.Since(startTime).Minutes())
		return results
	}

	log.Infof("Successful directory load from %v took %0.2f minutes",
		dataUrl, time.Since(startTime).Minutes())

	return nil
}

func loadZip(dataUrl string, companyData map[string]string, dataMap map[model.Ticker][]*model.Period) error {

	log.Infof("Starting load of zip archive %v...", dataUrl)
	startTime := time.Now()

	var results error

	// Open a zip archive for reading.
	r, err := zip.OpenReader(dataUrl)
	if err != nil {
		return errors.Wrap(err, "problem open zip archive")
	}
	defer func(r *zip.ReadCloser) {
		err := r.Close()
		if err != nil {
			results = multierror.Append(results, errors.Wrap(err, "problem closing zip reader"))
		}
	}(r)

	// Iterate through the files testInputData the archive,
	// printing some of their contents.
	for _, f := range r.File {
		names := strings.Split(f.Name, ".")
		rc, err := f.Open()
		if err != nil {
			return errors.Wrap(err, "problem open zip file")
		}

		reader := csv.NewReader(rc)
		err = loadData(strings.ToUpper(names[0]), reader, companyData, dataMap)
		if err != nil {
			results = multierror.Append(results, err)
		}
		err = rc.Close()
		if err != nil {
			results = multierror.Append(results, errors.Wrap(err, "problem closing zip file reader"))
		}
	}

	if results != nil {
		log.Infof("Completed zip archive load from %v with errors took %0.2f minutes",
			dataUrl, time.Since(startTime).Minutes())
		return results
	}

	log.Infof("Success loading zip archive from %v took %0.2f minutes",
		dataUrl, time.Since(startTime).Minutes())

	return nil
}

func loadFile(dataUrl string, companyData map[string]string, dataMap map[model.Ticker][]*model.Period) error {

	log.Infof("Starting load of file %v...", dataUrl)
	startTime := time.Now()

	_, file := filepath.Split(dataUrl)
	split := strings.Split(file, ".")

	ext := split[len(split)-1]
	if ext != "txt" && ext != "csv" {
		return fmt.Errorf("skipping unrecognized file extension of %v", ext)
	}

	// Skip this error and let the load return if the reader is invalid
	csvFile, _ := os.Open(dataUrl)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Errorf("unable to close file due to %v", err)
		}
	}(csvFile)

	reader := csv.NewReader(bufio.NewReader(csvFile))
	// split[0] should be the ticker symbol
	err := loadData(strings.ToUpper(split[0]), reader, companyData, dataMap)
	if err != nil {
		return errors.Wrapf(err, "Completed file load from %v with errors took %0.2f minutes",
			dataUrl, time.Since(startTime).Minutes())
	}

	log.Infof("Successful file load from %v took %0.2f minutes",
		dataUrl, time.Since(startTime).Minutes())

	return nil
}

func loadData(symbol string, r *csv.Reader, companyData map[string]string,
	dataMap map[model.Ticker][]*model.Period) error {

	var results error

	vals, err := r.ReadAll()
	if err != nil {
		return errors.Wrap(err, "error reading csv")
	}

	if vals == nil {
		return errors.New(fmt.Sprintf("empty or invalid CSV for '%v'", symbol))
	}

	var periods model.PeriodSlice
	var ticker model.Ticker

	ticker = model.Ticker{Symbol: symbol, Company: companyData[symbol]}
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

		p := model.Period{Symbol: symbol, Date: date, Open: open, High: high, Low: low, Close: cl, Volume: volume}
		periods = append(periods, &p)
	}

	sort.Sort(periods)

	dataMap[ticker] = periods

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
