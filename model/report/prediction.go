package report

type PredictionSeries struct {
	Name                string  `json:"name"`
	Pattern             string  `json:"priorPeriodsWere"`
	ProbabilityUp       float64 `json:"probabilityOfNextBeingUp"`
	ProbabilityDown     float64 `json:"probabilityOfNextBeingDown"`
	ProbabilityNoChange float64 `json:"probabilityOfNextBeingNoChange"`
}

type Prediction struct {
	TickerSymbol string              `json:"ticker"`
	FromDate     string              `json:"predictingFromDate"`
	NextDate     string              `json:"predictingDate"`
	Series       []*PredictionSeries `json:"series"`
}
