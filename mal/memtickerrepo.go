package mal

import (
	"fmt"
	"go-market-patterns/model"
	"go-market-patterns/model/report"
	"sync"
)

type MemTickerRepo struct {
	// map of ticker's symbol to a ticker pointer
	data  map[string]*model.Ticker
	mutex *sync.Mutex
}

func NewMemTickerRepo() *MemTickerRepo {
	return &MemTickerRepo{}
}

func (repo *MemTickerRepo) Init() {
	repo.data = make(map[string]*model.Ticker)
	repo.mutex = &sync.Mutex{}
}

func (repo *MemTickerRepo) CountAll() (int64, error) {
	return int64(len(repo.data)), nil
}

func (repo *MemTickerRepo) InsertOne(ticker *model.Ticker) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	repo.data[ticker.Symbol] = ticker
	return nil
}

func (repo *MemTickerRepo) InsertMany(data []*model.Ticker) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	for _, ticker := range data {
		repo.data[ticker.Symbol] = ticker
	}
	return nil
}

func (repo *MemTickerRepo) DropAndCreate() error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	repo.data = make(map[string]*model.Ticker)
	return nil
}

func (repo *MemTickerRepo) FindOne(symbol string) (*model.Ticker, error) {
	return repo.data[symbol], nil
}

func (repo *MemTickerRepo) FindOneCompanyNameBySymbol(symbol string) (string, error) {
	tmp := repo.data[symbol]
	if tmp == nil {
		return "", fmt.Errorf("company name not found for symbol '%v'", symbol)
	}

	return tmp.Company, nil
}

func (repo *MemTickerRepo) FindSymbols() []string {
	symbols := make([]string, len(repo.data))
	idx := 0
	for symbol, _ := range repo.data {
		symbols[idx] = symbol
		idx++
	}
	return symbols
}

func (repo *MemTickerRepo) FindSymbolsAndCompany() *report.TickerSymbolCompanySlice {
	slice := make(report.TickerSymbolCompanySlice, len(repo.data))
	idx := 0
	for _, ticker := range repo.data {
		slice[idx] = &report.TickerSymbolCompany{Symbol: ticker.Symbol, Company: ticker.Company}
	}
	return &slice
}
