package main

import (
	"github.com/hashicorp/go-multierror"
	"github.com/namsral/flag"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"market-patterns/config"
	"market-patterns/mal"
	"market-patterns/model"
	"market-patterns/tools"
	"os"
	"strings"
	"time"
)

var Repos *mal.Repos

func main() {

	var yamlConfig string
	flag.StringVar(&yamlConfig, "yaml-config", "", "YAML file containing configuration settings")
	flag.Parse()

	conf := config.Init(yamlConfig)
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
		err := truncAndLoad(conf.Options.DataFile, conf.Options.CompanyFile)
		if err != nil {
			log.Error(errors.Wrap(err, "problem truncating and loading"))
		}
	}

	if conf.Runtime.StartHttpServer {

		if conf.Runtime.HttpServerUrl == "" {
			log.Fatal("Invalid http-server-url")
		}
		// Start the profiler
		go startProfile()

		// Start the main api server
		start(conf)
	}
}

func truncAndLoad(dataFile, companyFile string) error {

	if dataFile == "" {
		log.Fatal("data-file must be specified for a trunc and load.")
	}

	if companyFile == "" {
		log.Fatal("company-file must be specified for a trunc and load.")
	}

	startTime := time.Now()

	log.Info("Deleting repos...")
	var dropErrors error
	err := Repos.PatternRepo.DeleteAll()
	if err != nil {
		dropErrors = multierror.Append(dropErrors, err)
	}
	err = Repos.PeriodRepo.DeleteAll()
	if err != nil {
		dropErrors = multierror.Append(dropErrors, err)
	}
	err = Repos.SeriesRepo.DeleteAll()
	if err != nil {
		dropErrors = multierror.Append(dropErrors, err)
	}
	err = Repos.TickerRepo.DeleteAll()
	if err != nil {
		dropErrors = multierror.Append(dropErrors, err)
	}
	if dropErrors != nil {
		return errors.Wrap(dropErrors, "unable to delete all repos")
	} else {
		log.Info("Success deleting repos.")
	}

	var loadErrors error
	dataMap := make(map[model.Ticker][]*model.Period)
	err = load(dataFile, companyFile, dataMap)
	if err != nil {
		loadErrors = multierror.Append(loadErrors, err)
	}

	err = train(3, dataMap)
	if err != nil {
		loadErrors = multierror.Append(loadErrors, err)
	}

	if loadErrors != nil {
		log.Infof("Completed trunc and load of %v with errors took %0.2f minutes",
			dataFile, time.Since(startTime).Minutes())
		log.Error(loadErrors)
	} else {
		log.Infof("Successful trunc and load of %v took %0.2f minutes",
			dataFile, time.Since(startTime).Minutes())
	}

	return nil
}
