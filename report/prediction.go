package report

type PredictionSeries struct {
	Name          string             `json:"name"`
	Pattern       string             `json:"priorPeriodsWere"`
	Probabilities map[string]float64 `json:"probabilityOfNextBeing"`
}

type Prediction struct {
	TickerSymbol string             `json:"ticker"`
	FromDate     string             `json:"predictingFromDate"`
	NextDate     string             `json:"predictingDate"`
	Series       []PredictionSeries `json:"series"`
}
