package main

import (
	log "github.com/sirupsen/logrus"
	"market-patterns/report"
	"market-patterns/utils"
	"net/http"
	_ "net/http/pprof"
)

func start() {
	fs := http.FileServer(http.Dir("static"))
	h := http.NewServeMux()
	h.HandleFunc("/api/predict", handlePredict)
	h.HandleFunc("/api/ticker-names", handlePredict)
	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
	log.Info("market-pattern server listening...")
	log.Fatal(http.ListenAndServe(":7666", h))
}

func handlePredict(w http.ResponseWriter, r *http.Request) {
	tickerNames := report.TickerNames{Names: Tickers.FindNames()}
	jsonData := utils.ToJsonBytes(tickerNames)
	r.Header.Set("Content-Type", "application/json")
	_, err := w.Write(jsonData)

	if err != nil {
		log.Errorf("unable to write response due to %v", err)
	}
}

func handleTickerNames(w http.ResponseWriter, r *http.Request) {

}

func startProfile() {
	log.Info("Starting profile server...")
	log.Fatal(http.ListenAndServe("localhost:6060", nil))
}
