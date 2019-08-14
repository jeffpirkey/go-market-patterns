package config

import (
	"github.com/namsral/flag"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type AppConfig struct {
	Runtime RuntimeConfig `yaml:"runtime"`
	Options OptionsConfig `yaml:"options"`
}

type OptionsConfig struct {
	StartHttpServer bool   `yaml:"start-http-server"`
	TruncLoad       bool   `yaml:"trunc-load"`
	DataFile        string `yaml:"data-file"`
	CompanyFile     string `yaml:"company-file"`
	PrintMDFile     string `yaml:"print-markdown"`
	Compute         int    `yaml:"compute"`
}

type RuntimeConfig struct {
	MongoDBUrl    string `yaml:"mongo-url"`
	MongoDBName   string `yaml:"mongo-dbname"`
	LogLevel      string `yaml:"log-level"`
	HttpServerUrl string `yaml:"http-server-url"`
}

func (c RuntimeConfig) Level() log.Level {
	level, err := log.ParseLevel(c.LogLevel)
	if err != nil {
		return 0
	}
	return level
}

var (
	initialized = false
	config      = &AppConfig{}
)

func Init(fileName string) *AppConfig {

	if initialized {
		return config
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("unable to load app configuration due to %v", err)
	}

	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "unable to process %v due to %v", fileName))
	}

	// Optional properties
	flag.BoolVar(&config.Options.StartHttpServer, "start-http-server", true,
		"start the http server, defaults to true")
	flag.BoolVar(&config.Options.TruncLoad, "trunc-load", false,
		"truncate and load, defaults to false")
	flag.StringVar(&config.Options.DataFile, "data-file", "",
		"load a csv, txt, zip file or load all files from a directory")
	flag.StringVar(&config.Options.CompanyFile, "company-file", "",
		"load symbol to company names")
	flag.StringVar(&config.Options.PrintMDFile, "print-markdown", "",
		"print markdown for the given symbol to output directory")
	flag.IntVar(&config.Options.Compute, "compute", 0,
		"compute the given series length, deleting the series if it exists")

	flag.Parse()

	if config.Runtime.Level() == log.DebugLevel {
		log.SetReportCaller(true)
	}
	log.SetLevel(config.Runtime.Level())

	if config.Options.Compute != 0 && (config.Options.Compute == 1 || config.Options.Compute < 0) {
		log.Fatal("compute length must be greater than 1")
	}

	initialized = true

	return config
}
