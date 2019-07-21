package graph

type StockSeries struct {
	Type  string          `json:"type"`
	Id    string          `json:"id"`
	Name  string          `json:"name"`
	Data  [][]interface{} `json:"data"`
	YAxis int             `json:"yAxis"`
}
