package config

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type AppConfig struct {
	Runtime RuntimeConfig `yaml:"runtime"`
}

type RuntimeConfig struct {
	MongoDBUrl  string `yaml:"mongo-db-url"`
	MongoDBName string `yaml:"mongo-db-name"`
}

func Init(fileName string) *AppConfig {

	//
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("unable to load app configuration due to %v", err)
	}
	config := AppConfig{}
	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Fatalf("unable to process config.yaml due to %v", err)
	}

	return &config
}
