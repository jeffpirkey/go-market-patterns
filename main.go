package main

import (
	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"
	"market-patterns/config"
	"market-patterns/mal"
)

/*
Requirements:
1) Lookup periods by ticker symbol
2) Sort periods by date
3) Build patterns for a ticker symbol using Up, NoChange, and Down
4) Count the number of Up, NoChange, and Down results for a pattern over time for a ticker symbol
5) Find a pattern for a ticker symbol
6) Find the most current period by date

*/

var Repos *mal.Repos

func main() {

	conf := config.Init("app-config.yaml")
	Repos = mal.New(conf)

	var load bool
	flag.BoolVar(&load, "load", false, "load and train")
	flag.Parse()

	if load {
		err := loadZip("data/stocks-small.zip")
		if err != nil {
			log.Fatal(err)
		}

		err = trainAll()
		if err != nil {
			log.Fatal(err)
		}
		err = trainAllSeries("3-period-series", "3 period series", 3)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("Completed load and train.")

		return
	}

	// Start the profiler
	go startProfile()

	// Start the main api server
	start()
}
