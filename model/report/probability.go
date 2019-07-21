package report

import "market-patterns/model"

type ProbabilityEdges struct {
	BestUpHigh       *model.Pattern
	BestDownHigh     *model.Pattern
	BestNoChangeHigh *model.Pattern
	BestUpLow        *model.Pattern
	BestDownLow      *model.Pattern
	BestNoChangeLow  *model.Pattern
}
