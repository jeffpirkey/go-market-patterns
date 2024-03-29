package mal

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-market-patterns/model"
	"sync"
)

type MemPatternRepo struct {
	data  map[string]map[string]map[int]*model.Pattern
	mutex *sync.Mutex
}

func NewMemPatternRepo() *MemPatternRepo {
	return &MemPatternRepo{}
}

func (repo *MemPatternRepo) Init() {
	repo.data = make(map[string]map[string]map[int]*model.Pattern)
	repo.mutex = &sync.Mutex{}
}

func (repo *MemPatternRepo) InsertMany(data []*model.Pattern) (int, error) {

	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	count := 0
	for _, pattern := range data {
		if symbolMap, found := repo.data[pattern.Symbol]; !found {
			repo.data[pattern.Symbol] = make(map[string]map[int]*model.Pattern)
			repo.data[pattern.Symbol][pattern.Value] = make(map[int]*model.Pattern)
			repo.data[pattern.Symbol][pattern.Value][pattern.Length] = pattern
			count++
		} else {
			if valueMap, found := symbolMap[pattern.Value]; !found {
				repo.data[pattern.Symbol][pattern.Value] = make(map[int]*model.Pattern)
				repo.data[pattern.Symbol][pattern.Value][pattern.Length] = pattern
				count++
			} else {
				if pattern, found := valueMap[pattern.Length]; !found {
					repo.data[pattern.Symbol][pattern.Value][pattern.Length] = pattern
					count++
				} else {
					log.Warnf("pattern already exists: %v", pattern)
				}
			}
		}
	}

	return count, nil
}

func (repo *MemPatternRepo) DeleteByLength(length int) error {
	panic("implement me")
}

func (repo *MemPatternRepo) DropAndCreate() error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	repo.data = make(map[string]map[string]map[int]*model.Pattern)
	return nil
}

func (repo *MemPatternRepo) FindBySymbol(symbol string) ([]*model.Pattern, error) {
	var patterns []*model.Pattern
	if symbolMap, found := repo.data[symbol]; found {
		for _, valueMap := range symbolMap {
			for _, pattern := range valueMap {
				patterns = append(patterns, pattern)
			}
		}
	} else {
		return patterns, fmt.Errorf("pattern not found for symbol '%v'", symbol)
	}

	return patterns, nil
}

func (repo *MemPatternRepo) FindOneBySymbolAndValueAndLength(symbol, value string, length int) (*model.Pattern, error) {

	if symbolMap, found := repo.data[symbol]; found {
		if valueMap, found := symbolMap[value]; found {
			if pattern, found := valueMap[length]; found {
				return pattern, nil
			}
		}
	}

	return nil,
		fmt.Errorf("pattern not found for symbol '%v', value '%v', and length '%v'", symbol, value, length)
}

func (repo *MemPatternRepo) FindHighestUpProbability(density model.PatternDensity) (*model.Pattern, error) {
	var max *model.Pattern
	for _, valueMap := range repo.data {
		for _, lengthMap := range valueMap {
			for _, pattern := range lengthMap {
				if max == nil {
					max = pattern
				} else {
					switch density {
					case model.PatternDensityLow:
						if pattern.UpCount/pattern.TotalCount > max.UpCount/max.TotalCount {
							max = pattern
						}
					case model.PatternDensityMedium:
						if pattern.TotalCount > 500 && pattern.UpCount/pattern.TotalCount > max.UpCount/max.TotalCount {
							max = pattern
						}
					case model.PatternDensityHigh:
						if pattern.TotalCount > 1000 && pattern.UpCount/pattern.TotalCount > max.UpCount/max.TotalCount {
							max = pattern
						}
					}
				}
			}
		}
	}

	return max, nil
}

func (repo *MemPatternRepo) FindHighestDownProbability(density model.PatternDensity) (*model.Pattern, error) {
	var max *model.Pattern
	for _, valueMap := range repo.data {
		for _, lengthMap := range valueMap {
			for _, pattern := range lengthMap {
				if max == nil {
					max = pattern
				} else {
					switch density {
					case model.PatternDensityLow:
						if pattern.UpCount/pattern.TotalCount > max.DownCount/max.TotalCount {
							max = pattern
						}
					case model.PatternDensityMedium:
						if pattern.TotalCount > 500 && pattern.DownCount/pattern.TotalCount > max.DownCount/max.TotalCount {
							max = pattern
						}
					case model.PatternDensityHigh:
						if pattern.TotalCount > 1000 && pattern.DownCount/pattern.TotalCount > max.DownCount/max.TotalCount {
							max = pattern
						}
					}
				}
			}
		}
	}

	return max, nil
}

func (repo *MemPatternRepo) FindHighestNoChangeProbability(density model.PatternDensity) (*model.Pattern, error) {
	var max *model.Pattern
	for _, valueMap := range repo.data {
		for _, lengthMap := range valueMap {
			for _, pattern := range lengthMap {
				if max == nil {
					max = pattern
				} else {
					switch density {
					case model.PatternDensityLow:
						if pattern.UpCount/pattern.TotalCount > max.NoChangeCount/max.TotalCount {
							max = pattern
						}
					case model.PatternDensityMedium:
						if pattern.TotalCount > 500 && pattern.NoChangeCount/pattern.TotalCount > max.NoChangeCount/max.TotalCount {
							max = pattern
						}
					case model.PatternDensityHigh:
						if pattern.TotalCount > 1000 && pattern.NoChangeCount/pattern.TotalCount > max.NoChangeCount/max.TotalCount {
							max = pattern
						}
					}
				}
			}
		}
	}

	return max, nil
}

func (repo *MemPatternRepo) FindLowestUpProbability(density model.PatternDensity) (*model.Pattern, error) {
	var max *model.Pattern
	for _, valueMap := range repo.data {
		for _, lengthMap := range valueMap {
			for _, pattern := range lengthMap {
				if max == nil {
					max = pattern
				} else {
					switch density {
					case model.PatternDensityLow:
						if pattern.UpCount/pattern.TotalCount < max.UpCount/max.TotalCount {
							max = pattern
						}
					case model.PatternDensityMedium:
						if pattern.TotalCount > 500 && pattern.UpCount/pattern.TotalCount < max.UpCount/max.TotalCount {
							max = pattern
						}
					case model.PatternDensityHigh:
						if pattern.TotalCount > 1000 && pattern.UpCount/pattern.TotalCount < max.UpCount/max.TotalCount {
							max = pattern
						}
					}
				}
			}
		}
	}

	return max, nil

}

func (repo *MemPatternRepo) FindLowestDownProbability(density model.PatternDensity) (*model.Pattern, error) {
	var max *model.Pattern
	for _, valueMap := range repo.data {
		for _, lengthMap := range valueMap {
			for _, pattern := range lengthMap {
				if max == nil {
					max = pattern
				} else {
					switch density {
					case model.PatternDensityLow:
						if pattern.UpCount/pattern.TotalCount < max.DownCount/max.TotalCount {
							max = pattern
						}
					case model.PatternDensityMedium:
						if pattern.TotalCount > 500 && pattern.DownCount/pattern.TotalCount < max.DownCount/max.TotalCount {
							max = pattern
						}
					case model.PatternDensityHigh:
						if pattern.TotalCount > 1000 && pattern.DownCount/pattern.TotalCount < max.DownCount/max.TotalCount {
							max = pattern
						}
					}
				}
			}
		}
	}

	return max, nil
}

func (repo *MemPatternRepo) FindLowestNoChangeProbability(density model.PatternDensity) (*model.Pattern, error) {
	var max *model.Pattern
	for _, valueMap := range repo.data {
		for _, lengthMap := range valueMap {
			for _, pattern := range lengthMap {
				if max == nil {
					max = pattern
				} else {
					switch density {
					case model.PatternDensityLow:
						if pattern.UpCount/pattern.TotalCount < max.NoChangeCount/max.TotalCount {
							max = pattern
						}
					case model.PatternDensityMedium:
						if pattern.TotalCount > 500 && pattern.NoChangeCount/pattern.TotalCount < max.NoChangeCount/max.TotalCount {
							max = pattern
						}
					case model.PatternDensityHigh:
						if pattern.TotalCount > 1000 && pattern.NoChangeCount/pattern.TotalCount < max.NoChangeCount/max.TotalCount {
							max = pattern
						}
					}
				}
			}
		}
	}

	return max, nil
}
