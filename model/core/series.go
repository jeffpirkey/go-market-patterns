package core

// This defines any series that have been run
type Series struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Length int    `json:"length"`
}
