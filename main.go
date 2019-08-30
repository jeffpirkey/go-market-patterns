package main

import (
	"errors"
	"github.com/hashicorp/go-multierror"
	log "github.com/sirupsen/logrus"
	"go-market-patterns/config"
	"go-market-patterns/mal"
	"go-market-patterns/tools"
	"os"
	"strings"
	"time"
)

var (
	Repos *mal.Repos
)

const (
	appConfig = "runtime-config.yaml"
)

func main() {

	conf := config.Init()

	Repos = mal.New(conf)

	if conf.Options.PrintMDFile != "" {

		err := os.MkdirAll("output", os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		fileName := "output" + string(os.PathSeparator) + conf.Options.PrintMDFile + ".md"
		log.Infof("Printing markdown for %v to %v", conf.Options.PrintMDFile, fileName)
		f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			log.Fatal(err)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				log.Error(err)
			}
		}(f)

		err = tools.PrintMarkdownPatterns(Repos, strings.ToUpper(conf.Options.PrintMDFile), f)
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("Successfully generated markdown for %v to %v", conf.Options.PrintMDFile, fileName)
	}

	if conf.Options.TruncLoad {
		log.Info("Started trunc and load")
		startTime := time.Now()
		err := truncAndLoad(conf.Options.DataFile, conf.Options.CompanyFile, conf.Options.Compute)
		if err != nil {
			log.Error(err)
		}
		log.Infof("Completed trunc and load took %0.2f minutes",
			time.Since(startTime).Minutes())
	} else if conf.Options.Compute > 1 {
		log.Infof("Started computing all periods with length %v...", conf.Options.Compute)
		startTime := time.Now()
		err := truncAndComputeAllSeries(conf.Options.Compute)
		if err != nil {
			log.WithError(err).Errorf(
				"Completed training all periods with length %v with errors took %0.2f minutes",
				conf.Options.Compute, time.Since(startTime).Minutes())
		} else {
			log.Infof("Completed training all periods with length %v took %0.2f minutes",
				conf.Options.Compute, time.Since(startTime).Minutes())
		}
	}

	if conf.Options.StartHttpServer {

		if conf.Runtime.HttpServerUrl == "" {
			log.Fatal("Invalid http-server-url")
		}
		// Start the profiler
		go startProfile()

		// Start the main api server
		start(conf)
	}
}

// This function deletes all the repo data, reloads from the given data and company files.
// After loading, the one-day train is executed against the periods.
func truncAndLoad(dataFile, companyFile string, computeLength int) error {

	if dataFile == "" {
		log.Fatal("data-file must be specified for a trunc and load.")
	}

	if companyFile == "" {
		log.Fatal("company-file must be specified for a trunc and load.")
	}

	log.Info("Dropping and recreating repos...")
	startTime := time.Now()
	var dropErrors error
	if err := Repos.PatternRepo.DropAndCreate(); err != nil {
		dropErrors = multierror.Append(dropErrors, err)
	}
	if err := Repos.PeriodRepo.DropAndCreate(); err != nil {
		dropErrors = multierror.Append(dropErrors, err)
	}
	if err := Repos.SeriesRepo.DropAndCreate(); err != nil {
		dropErrors = multierror.Append(dropErrors, err)
	}
	if err := Repos.TickerRepo.DropAndCreate(); err != nil {
		dropErrors = multierror.Append(dropErrors, err)
	}

	if dropErrors != nil {
		log.WithError(dropErrors).Errorf("Completed recreating repos with errors took %0.2f minutes",
			time.Since(startTime).Minutes())
		return errors.New("unable to continue trunc and load due to repo recreate issues")
	}

	log.Infof("Completed recreating repos took %0.2f minutes",
		time.Since(startTime).Minutes())

	return load(dataFile, companyFile, computeLength)
}
