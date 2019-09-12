package report

import (
	"go-market-patterns/model/core"
)

type ProbabilityEdges struct {
	BestUpHigh       *core.Pattern `json:"bestUpHigh"`
	BestDownHigh     *core.Pattern `json:"bestDownHigh"`
	BestNoChangeHigh *core.Pattern `json:"bestNoChangeHigh"`
	BestUpLow        *core.Pattern `json:"bestUpLow"`
	BestDownLow      *core.Pattern `json:"bestDownLow"`
	BestNoChangeLow  *core.Pattern `json:"bestNoChangeLow"`
}
