package model

const (
	Up         = "U"
	NoChange   = "N"
	Down       = "D"
	NotDefined = "x"
)

// For a given sequence of Result types, a pattern records the number of times that the next sequence is Up, Down, or
// NoChange.
type Pattern struct {
	Results map[string]int
}

// *********************************************************
//   Pattern methods
// *********************************************************

func (p *Pattern) FindAll() map[string]int {
	if p.Results == nil {
		p.Results = make(map[string]int)
	}
	return p.Results
}

func (p *Pattern) Find(r string) int {
	if p.Results == nil {
		p.Results = make(map[string]int)
	}
	return p.Results[r]
}

func (p *Pattern) TotalCount() int {
	sum := 0
	for _, v := range p.Results {
		sum += v
	}
	return sum
}

func (p *Pattern) Inc(result string) {
	if p.Results == nil {
		p.Results = make(map[string]int)
	}
	p.Results[result]++
}

// This function calculates a Result using 'prev' as the previous value compared to the current date in 'cur'
func Calc(prev, cur float64) string {
	if prev < cur {
		return Up
	} else if prev > cur {
		return Down
	}
	return NoChange
}
