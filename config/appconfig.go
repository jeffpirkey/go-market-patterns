package config

import (
	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type AppConfig struct {
	Runtime RuntimeConfig `yaml:"runtime"`
	Options OptionsConfig `yaml:"options"`
}

type OptionsConfig struct {
	TruncLoad   bool   `yaml:"trunc-load"`
	DataFile    string `yaml:"data-file"`
	CompanyFile string `yaml:"company-file"`
	PrintMDFile string `yaml:"print-markdown"`
}

type RuntimeConfig struct {
	MongoDBUrl      string `yaml:"mongo-url"`
	MongoDBName     string `yaml:"mongo-dbname"`
	LogLevel        string `yaml:"log-level"`
	StartHttpServer bool   `yaml:"start-http-server"`
}

func (c RuntimeConfig) Level() log.Level {
	level, err := log.ParseLevel(c.LogLevel)
	if err != nil {
		return 0
	}
	return level
}

func Init(yamlFileName string) *AppConfig {

	config := AppConfig{}

	if yamlFileName == "" {
		flag.BoolVar(&config.Runtime.StartHttpServer, "start-http-server", true,
			"start http server, defaults to true")
		flag.BoolVar(&config.Options.TruncLoad, "trunc-load", false,
			"truncate and load, defaults to false")
		flag.StringVar(&config.Runtime.LogLevel, "log-level", "INFO",
			"set logging level, defaults to 'INFO'")
		flag.StringVar(&config.Options.CompanyFile, "data-file", "",
			"load a csv, txt, zip file or load all files from a directory")
		flag.StringVar(&config.Options.DataFile, "company-file", "",
			"load symbol to company names")
		flag.StringVar(&config.Options.PrintMDFile, "print-markdown", "",
			"print markdown for the given symbol to output directory")
		flag.StringVar(&config.Runtime.MongoDBName, "mongo-dbname", "",
			"mongodb database name")
		flag.StringVar(&config.Runtime.MongoDBUrl, "mongo-url", "",
			"mongodb url, such as 'mongodb://localhost:27017'")
		flag.Parse()
	} else {
		data, err := ioutil.ReadFile(yamlFileName)
		if err != nil {
			log.Fatalf("unable to load app configuration due to %v", err)
		}

		err = yaml.Unmarshal([]byte(data), &config)
		if err != nil {
			log.Fatalf("unable to process config.yaml due to %v", err)
		}
	}

	if config.Runtime.Level() == log.DebugLevel {
		log.SetReportCaller(true)
	}
	log.SetLevel(config.Runtime.Level())

	return &config
}
