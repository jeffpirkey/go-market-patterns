package config

import (
	"errors"
	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"
	"strconv"
)

type AppConfig struct {
	Runtime RuntimeConfig `yaml:"runtime"`
	Options OptionsConfig `yaml:"options"`
}

type OptionsConfig struct {
	StartHttpServer bool       `yaml:"start-http-server"`
	TruncLoad       bool       `yaml:"trunc-load"`
	DataFile        string     `yaml:"data-file"`
	CompanyFile     string     `yaml:"company-file"`
	PrintMDFile     string     `yaml:"print-markdown"`
	Compute         arrayFlags `yaml:"compute"`
}

type RuntimeConfig struct {
	DbConnect     string `yaml:"db-connect"`
	MongoDbName   string `yaml:"mongo-db-name"`
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

type arrayFlags []int

func (i *arrayFlags) String() string {
	return "array flags"
}

func (i *arrayFlags) Set(value string) error {
	tmp, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	if tmp != 0 && (tmp == 1 || tmp < 0) {
		return errors.New("compute length must be greater than 1")
	}

	*i = append(*i, tmp)
	return nil
}

var (
	initialized    = false
	config         = &AppConfig{}
	computeLengths arrayFlags
)

func Init() *AppConfig {

	if initialized {
		return config
	}

	// Runtime properties
	flag.StringVar(&config.Runtime.DbConnect, "db-connect", "memory",
		"the db connection protocol, defaults to memory")
	flag.StringVar(&config.Runtime.MongoDbName, "mongo-db-name", "marketPatterns",
		"the database name to use in mongo, defaults to marketPatterns")
	flag.StringVar(&config.Runtime.LogLevel, "log-level", "DEBUG",
		"the logging level, defaults to DEBUG")
	flag.StringVar(&config.Runtime.HttpServerUrl, "http-server-url", ":8081",
		"the http server url, defaults to :8081")
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
	flag.Var(&config.Options.Compute, "compute",
		"compute the given series length, deleting the series if it exists")

	flag.Parse()

	if config.Runtime.Level() == log.DebugLevel {
		log.SetReportCaller(true)
	}
	log.SetLevel(config.Runtime.Level())

	initialized = true

	return config
}
