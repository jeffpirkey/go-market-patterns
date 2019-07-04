package report

type PredictionSeries struct {
	Pattern       string             `json:"prior-periods-were"`
	Probabilities map[string]float64 `json:"probability-of-next-being"`
}

type Prediction struct {
	TickerSymbol string                      `json:"ticker"`
	FromDate     string                      `json:"predicting-from-date"`
	NextDate     string                      `json:"predicting-date"`
	Series       map[string]PredictionSeries `json:"series"`
}
