package model

type Result int

const (
	Up Result = iota
	NoChange
	Down
	NotDefined
)

// For a given sequence of Result types, a pattern records the number of times that the next sequence is Up, Down, or
// NoChange.
type Pattern struct {
	results map[Result]int
}

// *********************************************************
//   Pattern methods
// *********************************************************

func (p *Pattern) FindAll() map[Result]int {
	if p.results == nil {
		p.results = make(map[Result]int)
	}
	return p.results
}

func (p *Pattern) Find(r Result) int {
	if p.results == nil {
		p.results = make(map[Result]int)
	}
	return p.results[r]
}

func (p *Pattern) TotalCount() int {
	sum := 0
	for _, v := range p.results {
		sum += v
	}
	return sum
}

func (p *Pattern) Inc(result Result) {
	if p.results == nil {
		p.results = make(map[Result]int)
	}
	p.results[result]++
}

func (r Result) String() string {
	return [...]string{"U", "N", "D", "x"}[r]
}

// This function calculates a Result using 'prev' as the previous value compared to the current date in 'cur'
func Calc(prev, cur float64) Result {
	if prev < cur {
		return Up
	} else if prev > cur {
		return Down
	}
	return NoChange
}
