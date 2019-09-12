package mal

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-market-patterns/config"
	"go-market-patterns/model/core"
	"go-market-patterns/model/report"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"strings"
	"testing"
)

type SortDirection int

const (
	SortDsc SortDirection = iota
	SortAsc
)

var (
	client     *mongo.Client
	UpdateOpt  = options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	ReplaceOpt = options.FindOneAndReplace().SetUpsert(true).SetReturnDocument(options.After)
)

// Exported type for repository access
type Repos struct {
	client          *mongo.Client
	config          *config.AppConfig
	TickerRepo      TickerRepo
	PatternRepo     PatternRepo
	PeriodRepo      PeriodRepo
	SeriesRepo      SeriesRepo
	GraphController *GraphController
}

func New(config *config.AppConfig) *Repos {
	r := Repos{}
	r.Init(config)
	return &r
}

func (repos *Repos) Init(config *config.AppConfig) {

	repos.config = config

	if strings.HasPrefix(config.Runtime.DbConnect, "memory") {
		repos.TickerRepo = NewMemTickerRepo()
		repos.SeriesRepo = NewMemSeriesRepo()
		repos.PatternRepo = NewMemPatternRepo()
		repos.PeriodRepo = NewMemPeriodRepo()
	} else if strings.HasPrefix(config.Runtime.DbConnect, "mongodb") {
		var err error
		client, err = mongo.NewClient(options.Client().ApplyURI(config.Runtime.DbConnect))
		if err != nil {
			log.Fatalf("unable to create mongodb client due to %v", err)
		}

		err = client.Connect(context.TODO())
		if err != nil {
			log.Fatalf("unable to connect to mongodb at %v due to %v", config.Runtime.DbConnect, err)
		}

		err = client.Ping(context.TODO(), readpref.Primary())
		if err != nil {
			log.Fatalf("unable to ping mongodb at %v due to %v", config.Runtime.DbConnect, err)
		}

		coll := client.Database(config.Runtime.MongoDbName).Collection("tickers")
		repos.TickerRepo = NewMongoTickerRepo(coll)

		coll = client.Database(config.Runtime.MongoDbName).Collection("patterns")
		repos.PatternRepo = NewMongoPatternRepo(coll)

		coll = client.Database(config.Runtime.MongoDbName).Collection("periods")
		repos.PeriodRepo = NewMongoPeriodRepo(coll)

		coll = client.Database(config.Runtime.MongoDbName).Collection("series")
		repos.SeriesRepo = NewMongoSeriesRepo(coll)

	} else {
		log.Fatalf("unrecognized db protocol '%v'", config.Runtime.DbConnect)
	}

	repos.TickerRepo.Init()
	repos.PatternRepo.Init()
	repos.PeriodRepo.Init()
	repos.SeriesRepo.Init()

	repos.GraphController = &GraphController{repos.PeriodRepo, repos.PatternRepo}

}

func (repos *Repos) DropAll(t *testing.T) {

	if t == nil {
		log.Error("only able to delete databases when in test mode")
		return
	}

	if strings.HasPrefix(repos.config.Runtime.DbConnect, "memory") {
		err := repos.PatternRepo.DropAndCreate()
		if err != nil {
			log.Errorf("problem dropping pattern repo: %v", err)
		}
		err = repos.PeriodRepo.DropAndCreate()
		if err != nil {
			log.Errorf("problem dropping pattern repo: %v", err)
		}
		err = repos.SeriesRepo.DropAndCreate()
		if err != nil {
			log.Errorf("problem dropping pattern repo: %v", err)
		}
		err = repos.TickerRepo.DropAndCreate()
		if err != nil {
			log.Errorf("problem dropping pattern repo: %v", err)
		}
	} else if strings.HasPrefix(repos.config.Runtime.DbConnect, "mongo") {
		db := client.Database(repos.config.Runtime.DbConnect)
		err := db.Drop(context.TODO())
		if err != nil {
			log.Errorf("unable to drop database %v due to %v", repos.config.Runtime.DbConnect, err)
		}
	}
}

func CreateCollection(c *mongo.Collection, doc interface{}) (bool, error) {

	log.Infof("Checking if collection '%v' needs to be created", c.Name())
	tmp := c.FindOne(context.TODO(), bson.D{})
	created := false
	if tmp == nil {
		insertResult, err := c.InsertOne(context.TODO(), doc)
		if err != nil {
			return false, err
		}
		delResult, err := c.DeleteOne(context.TODO(), bson.D{{"_id", insertResult.InsertedID}})
		if err != nil {
			return false, err
		}
		if delResult.DeletedCount == 0 {
			return false, fmt.Errorf("problem creating collection '%v' with priming read", c.Name())
		}
		created = true
		log.Infof("Created collection '%v'", c.Name())
	} else {
		log.Infof("Collection '%v' already exists", c.Name())
	}

	return created, nil
}

type PatternRepo interface {
	Init()
	InsertMany(data []*core.Pattern) (int, error)
	DeleteByLength(length int) error
	DropAndCreate() error
	FindBySymbol(symbol string) ([]*core.Pattern, error)
	FindBySymbolAndLength(symbol string, length int) ([]*core.Pattern, error)
	FindOneBySymbolAndValueAndLength(symbol, value string, length int) (*core.Pattern, error)
	FindHighestUpProbability(density core.PatternDensity) (*core.Pattern, error)
	FindHighestDownProbability(density core.PatternDensity) (*core.Pattern, error)
	FindHighestNoChangeProbability(density core.PatternDensity) (*core.Pattern, error)
	FindLowestUpProbability(density core.PatternDensity) (*core.Pattern, error)
	FindLowestDownProbability(density core.PatternDensity) (*core.Pattern, error)
	FindLowestNoChangeProbability(density core.PatternDensity) (*core.Pattern, error)
}

type PeriodRepo interface {
	Init()
	InsertMany(data []*core.Period) (int, error)
	DropAndCreate() error
	FindBySymbol(symbol string, sort SortDirection) (core.PeriodSlice, error)
}

type TickerRepo interface {
	Init()
	CountAll() (int64, error)
	InsertOne(ticker *core.Ticker) error
	InsertMany(data []*core.Ticker) error
	DropAndCreate() error
	FindOne(symbol string) (*core.Ticker, error)
	FindOneCompanyNameBySymbol(symbol string) (string, error)
	FindSymbols() []string
	FindSymbolCompanySliceSortAsc() *report.TickerSymbolCompanySlice
}

type SeriesRepo interface {
	Init()
	FindBySymbol(symbol string) ([]*core.Series, error)
	FindOneBySymbolAndLength(symbol string, length int) (*core.Series, error)
	FindNameLengthSliceBySymbol(symbol string) *report.SeriesNameLengthSlice
	InsertOne(data *core.Series) error
	DeleteOne(data *core.Series) error
	DeleteByLength(length int) error
	DropAndCreate() error
}
