package graph

type PatternDensity struct {
	Categories []string `json:"categories"`
	Totals     []int    `json:"totals"`
	Ups        []int    `json:"ups"`
	Downs      []int    `json:"downs"`
	NoChanges  []int    `json:"nochanges"`
}
