package mal

import (
	log "github.com/sirupsen/logrus"
	"go-market-patterns/model"
	"sort"
	"sync"
	"time"
)

type MemPeriodRepo struct {
	data  map[string]map[time.Time]*model.Period
	mutex *sync.Mutex
}

func NewMemPeriodRepo() *MemPeriodRepo {
	return &MemPeriodRepo{}
}

func (repo *MemPeriodRepo) Init() {
	repo.data = make(map[string]map[time.Time]*model.Period)
	repo.mutex = &sync.Mutex{}
}

func (repo *MemPeriodRepo) InsertMany(data []*model.Period) (int, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	count := 0
	for _, period := range data {
		if symbolMap, found := repo.data[period.Symbol]; found {
			// Symbol exists, check if period exists for the given time
			if _, found := symbolMap[period.Date]; !found {
				// Time not found, so insert
				repo.data[period.Symbol][period.Date] = period
				count++
			} else {
				log.Errorf("period already exists: %v", period)
			}
		} else {
			// Symbol not in map, so create it
			repo.data[period.Symbol] = make(map[time.Time]*model.Period)
			repo.data[period.Symbol][period.Date] = period
			count++
		}
	}

	return count, nil
}

func (repo *MemPeriodRepo) DropAndCreate() error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	repo.data = make(map[string]map[time.Time]*model.Period)
	return nil
}

func (repo *MemPeriodRepo) FindBySymbol(symbol string, sortDir SortDirection) (model.PeriodSlice, error) {

	var tmp model.PeriodSlice
	if symbolMap, found := repo.data[symbol]; found {
		for _, period := range symbolMap {
			tmp = append(tmp, period)
		}
	}

	if sortDir == SortDsc {
		sort.Reverse(tmp)
	} else {
		sort.Sort(tmp)
	}

	return tmp, nil
}
