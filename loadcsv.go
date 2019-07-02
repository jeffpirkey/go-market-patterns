package main

import (
	"encoding/csv"
	"log"
	"market-patterns/model"
	"sort"
	"strconv"
	"time"
)

const (
	timeFormat = "2006-01-02"
)

func load(name string, r *csv.Reader) {

	vals, err := r.ReadAll()
	if err != nil {
		log.Fatalf("error reading csv due to %v", err)
	}

	for i, v := range vals {

		if i == 0 {
			// skip header line
			continue
		}

		pv, err := strconv.ParseFloat(v[1], 64)
		if err != nil {
			log.Fatalf("error parsing csv value due ot %v", err)
		}
		pt, err := time.Parse(timeFormat, v[0])
		if err != nil {
			log.Fatalf("error parsing csv time due ot %v", err)
		}
		p := model.Period{pv, pt}

		Periods[name] = append(Periods[name], p)
	}

	sort.Sort(Periods[name])

}
