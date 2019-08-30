package main

import (
	"archive/zip"
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/rcrowley/go-metrics"
	log "github.com/sirupsen/logrus"
	"go-market-patterns/model"
	"go-market-patterns/utils"
	"golang.org/x/sync/semaphore"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type CsvField int

const (
	timeFormat = "2006-01-02"
)

const (
	csvDate CsvField = iota
	csvOpen
	csvHigh
	csvLow
	csvClose
	csvVolume
)

type CsvRow struct {
	Date   time.Time
	High   float64
	Low    float64
	Open   float64
	Close  float64
	Volume int
}

var (
	loadRegistry = metrics.NewRegistry()
	loadMutex    = &sync.Mutex{}
)

// load attempts to read, persist, and optionally pattern compute
// ticker data from the given url.
// The companyFile parameter is used to correlate ticker symbol to
// an actual company name.
// If a computeLength of greater than 1 is specified, then the
// compute patterns are run for that length.
// If an error is returned, this indicates that some portion
// of the process failed, but maybe not all.
func load(url, companyFile string, computeLengths []int) error {

	// Allow only one load to be happening at a time
	loadMutex.Lock()
	defer loadMutex.Unlock()

	loadRegistry.UnregisterAll()

	startTime := time.Now()
	log.Infof("Started load of company data from %v", companyFile)
	companyData, err := loadCompanies(companyFile)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"duration": time.Since(startTime).Seconds(),
			"company-file": companyFile}).Error("Completed company data load")
		return errors.Wrapf(err, "loading company data")
	}

	startTime = time.Now()
	fi, err := os.Stat(url)
	if err != nil {
		return errors.Wrapf(err, "checking url type")
	}

	var loadErrors error
	baseLogger := log.WithField("url", url)
	switch mode := fi.Mode(); {
	case mode.IsDir():
		baseLogger := baseLogger.WithField("url-type", "directory")
		baseLogger.Infof("Started load")
		loadErrors = loadDir(url, companyData, computeLengths)

	case mode.IsRegular():
		if utils.IsZip(url) {
			baseLogger := baseLogger.WithField("url-type", "archive")
			baseLogger.Infof("Started load")
			loadErrors = loadZip(url, companyData, computeLengths)
		} else {
			baseLogger := baseLogger.WithField("url-type", "file")
			baseLogger.Infof("Started load")
			loadErrors = loadFile(url, companyData, computeLengths)
		}
	default:
		err = errors.New("unrecognized load type")
	}

	return loadErrors
}

func loadCompanies(fileName string) (map[string]string, error) {

	var data map[string]string

	csvFile, err := os.Open(fileName)
	if err != nil {
		return data, errors.Wrap(err, "problem loading company data ")
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
		sym := strings.ReplaceAll(v[0], "-", "")
		data[sym] = v[1]
	}

	return data, nil
}

func loadDir(dataUrl string, companyData map[string]string, computeLengths []int) error {

	files, err := ioutil.ReadDir(dataUrl)
	if err != nil {
		return errors.Wrapf(err, "unable to load directory %v", dataUrl)
	}

	var loadErrors error
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
		symbol := strings.ToUpper(split[0])
		symbol = strings.ReplaceAll(symbol, "-", ".")
		symbol = strings.ReplaceAll(symbol, "_", "-")
		companyName := companyData[symbol]
		err = loadAndTrainData(symbol, companyName, reader, computeLengths)
		if err != nil {
			loadErrors = multierror.Append(loadErrors, err)
		}
		// Not a critical error, so don't add
		err = csvFile.Close()
		if err != nil {
			log.WithError(err).Error("unable to close csv file")
		}
	}

	return loadErrors
}

func loadZip(dataUrl string, companyData map[string]string, computeLengths []int) error {

	// Open a zip archive for reading.
	r, err := zip.OpenReader(dataUrl)
	if err != nil {
		return errors.Wrap(err, "problem open zip archive")
	}
	defer func(r *zip.ReadCloser) {
		err := r.Close()
		if err != nil {
			log.WithError(err).Errorf("problem closing zip reader")
		}
	}(r)

	var loadErrors error
	// Setup work pool
	ctx := context.TODO()
	var (
		maxWorkers = runtime.GOMAXPROCS(0)
		sem        = semaphore.NewWeighted(int64(maxWorkers))
	)

	for _, zipFile := range r.File {

		if err := sem.Acquire(ctx, 1); err != nil {
			return errors.New("failed to acquire semaphore during lock")
		}

		go func(zf *zip.File) {
			defer sem.Release(1)

			rc, err := zf.Open()
			if err != nil {
				loadErrors = multierror.Append(loadErrors,
					errors.Wrapf(err, "failed opening file in zip archive '%v'", zf.Name))
			}
			defer handleClose(rc)

			names := strings.Split(zf.Name, ".")
			// split[0] should be the ticker symbol
			symbol := strings.ToUpper(names[0])
			symbol = strings.ReplaceAll(symbol, "-", ".")
			symbol = strings.ReplaceAll(symbol, "_", "-")
			companyName := companyData[symbol]
			reader := csv.NewReader(rc)
			err = loadAndTrainData(symbol, companyName, reader, computeLengths)
			if err != nil {
				loadErrors = multierror.Append(loadErrors, err)
			}
		}(zipFile)
	}

	if err := sem.Acquire(ctx, int64(maxWorkers)); err != nil {
		log.WithError(err).Error("failed to acquire semaphore during unlock")
	}

	return loadErrors
}

func loadFile(dataUrl string, companyData map[string]string, computeLengths []int) error {

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
	symbol := strings.ToUpper(split[0])
	symbol = strings.ReplaceAll(symbol, "-", ".")
	symbol = strings.ReplaceAll(symbol, "_", "-")
	companyName := companyData[symbol]
	return loadAndTrainData(symbol, companyName, reader, computeLengths)
}

// loadAndTrainData parses the csv data from the given reader.
// After parsing, the periods are trained using the day-to-day algorithm.
// Optionally, if a computeLength greater than 1 is provide, the patterns are
// computed for the periods.
// loadAndTrainData is successful when the related Ticker is persisted to the repo;
// otherwise an error is returned.
func loadAndTrainData(symbol, companyName string, r *csv.Reader,
	computeLengths []int) error {

	vals, err := r.ReadAll()
	if err != nil {
		return errors.Wrap(err, "error reading csv")
	}

	if vals == nil {
		return errors.New(fmt.Sprintf("empty or invalid CSV for '%v'", symbol))
	}

	// Only allow if there are enough periods for train and compute
	if len(vals) <= 2 {
		return errors.New(fmt.Sprintf("[%v] not enough periods", symbol))
	}

	var periods model.PeriodSlice
	for i, v := range vals {

		var parseErrors error

		if i == 0 {
			// TODO jpirkey build header to field map index here
			continue
		}

		row := CsvRow{}
		if row.Date, err = convertTime(v[csvDate]); err != nil {
			parseErrors = multierror.Append(parseErrors, errors.Wrapf(err, "[%v] date field", symbol))
		}
		if row.Open, err = convertFloat(v[csvOpen]); err != nil {
			parseErrors = multierror.Append(parseErrors, errors.Wrapf(err, "[%v] open field", symbol))
		}
		if row.High, err = convertFloat(v[csvHigh]); err != nil {
			parseErrors = multierror.Append(parseErrors, errors.Wrapf(err, "[%v] high field", symbol))
		}
		if row.Low, err = convertFloat(v[csvLow]); err != nil {
			parseErrors = multierror.Append(parseErrors, errors.Wrapf(err, "[%v] low field", symbol))
		}
		if row.Close, err = convertFloat(v[csvClose]); err != nil {
			parseErrors = multierror.Append(parseErrors, errors.Wrapf(err, "[%v] close field", symbol))
		}
		if row.Volume, err = convertInt(v[csvVolume]); err != nil {
			parseErrors = multierror.Append(parseErrors, errors.Wrapf(err, "[%v] volume field", symbol))
		}

		if parseErrors == nil {
			p := model.Period{Symbol: symbol, Date: row.Date, Open: row.Open, High: row.High,
				Low: row.Low, Close: row.Close, Volume: row.Volume}
			periods = append(periods, &p)
		} else {
			log.Warn(parseErrors)
		}
	}

	if len(periods) < 2 {
		return errors.New(fmt.Sprintf("[%v] not enough parsed periods", symbol))
	}

	sort.Sort(periods)

	timer := metrics.GetOrRegisterTimer("training-timer", loadRegistry)
	timer.Time(func() { trainDaily(periods) })

	insertCount, err := Repos.PeriodRepo.InsertMany(periods)
	if err != nil {
		return errors.Wrapf(err, "[%v] inserting periods", symbol)
	}
	if len(periods) != insertCount {
		return fmt.Errorf("[%v] periods parsed count does not match inserted count", symbol)
	}

	ticker := model.Ticker{Symbol: symbol, Company: companyName}
	err = Repos.TickerRepo.InsertOne(&ticker)
	if err != nil {
		return errors.Wrapf(err, "[%v] inserting ticker", symbol)
	}

	if len(computeLengths) > 0 {
		for _, computeLength := range computeLengths {
			if len(periods) < computeLength+1 {
				log.Warnf("[%v] not enough periods to compute %v length series",
					symbol, computeLength)
			} else {
				timer := metrics.GetOrRegisterTimer("compute-timer", loadRegistry)
				timer.Time(func() { computeSeries(computeLength, symbol, periods) })
			}
		}
	}

	return nil
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

func handleClose(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.WithError(err).Error("failed closing")
	}
}
