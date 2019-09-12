package report

type Series struct {
	Series *SeriesNameLengthSlice `json:"series"`
}

type SeriesNameLength struct {
	Name   string `json:"name"`
	Length int    `json:"length"`
}

// Slice of Periods used for sorting and other access methods
type SeriesNameLengthSlice []*SeriesNameLength

// *********************************************************
//   PeriodSlice methods
// *********************************************************

func (p SeriesNameLengthSlice) Len() int {
	return len(p)
}

func (p SeriesNameLengthSlice) Less(i, j int) bool {
	return p[i].Length < p[j].Length
}

func (p SeriesNameLengthSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p SeriesNameLengthSlice) Last() *SeriesNameLength {
	return p[len(p)-1]
}

// Returns pointers to the items from the end of the slice.
func (p SeriesNameLengthSlice) LastByRange(l int) []*SeriesNameLength {
	return p[(len(p) - l):]
}
