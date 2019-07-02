package model

type Result string

const (
	Up       Result = "U"
	NoChange Result = "N"
	Down     Result = "D"
)

type pattern struct {
	results map[Result]int
}

type ticker struct {
	patterns map[string]pattern
}

type Tickers struct {
	tickers map[string]ticker
}

func NewTickers() Tickers {
	return Tickers{make(map[string]ticker)}
}

func (t *Tickers) Find(tsym string) ticker {

	x, ok := t.tickers[tsym]
	if ok {
		return x
	}
	t.tickers[tsym] = ticker{make(map[string]pattern)}

	return t.tickers[tsym]
}

func (t *ticker) FindAll() map[string]pattern {
	return t.patterns
}

func (t *ticker) Find(name string) pattern {
	x, ok := t.patterns[name]
	if ok {
		return x
	}

	t.patterns[name] = pattern{make(map[Result]int)}

	return t.patterns[name]
}

func (p *pattern) FindAll() map[Result]int {
	return p.results
}

func (p *pattern) Find(r Result) int {
	return p.results[r]
}

func (p *pattern) TotalCount() int {
	sum := 0

	for _, v := range p.results {
		sum += v
	}

	return sum
}

func (p *pattern) Inc(result Result) {

	p.results[result]++
}

// This function calculates the PeriodResult using p2 as the current timeSlice compared to the previous date in p1
func Calc(prev, cur float64) Result {

	if prev < cur {
		return Up
	} else if prev > cur {
		return Down
	}

	return NoChange
}
