package mal

import (
	"fmt"
	"go-market-patterns/model"
	"sync"
)

type MemSeriesRepo struct {
	// map of series' symbol to a series pointer
	data  map[string][]*model.Series
	mutex *sync.Mutex
}

func NewMemSeriesRepo() *MemSeriesRepo {
	return &MemSeriesRepo{}
}

func (repo *MemSeriesRepo) Init() {
	repo.data = make(map[string][]*model.Series)
	repo.mutex = &sync.Mutex{}
}

func (repo *MemSeriesRepo) FindBySymbol(symbol string) ([]*model.Series, error) {
	return repo.data[symbol], nil
}

func (repo *MemSeriesRepo) InsertOne(data *model.Series) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	for _, series := range repo.data[data.Symbol] {
		if series.Length == data.Length {
			return fmt.Errorf("series with length %v already exists", series.Length)
		}
	}

	repo.data[data.Symbol] = append(repo.data[data.Symbol], data)

	return nil
}

func (repo *MemSeriesRepo) DeleteOne(data *model.Series) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	seriesAry := repo.data[data.Symbol]
	removeIdx := -1
	for idx, series := range seriesAry {
		if series.Length == data.Length {
			removeIdx = idx
			break
		}
	}

	if removeIdx >= 0 {
		seriesAry[removeIdx] = seriesAry[len(seriesAry)-1]
	}

	return nil
}

func (repo *MemSeriesRepo) DeleteByLength(length int) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	for _, ary := range repo.data {
		for idx, series := range ary {
			if series.Length == length {
				ary[idx] = ary[len(ary)-1]
				break
			}
		}
	}

	return nil
}

func (repo *MemSeriesRepo) DropAndCreate() error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	repo.data = make(map[string][]*model.Series)
	return nil
}
