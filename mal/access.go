package mal

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"market-patterns/config"
	"market-patterns/model"
	"market-patterns/model/report"
	"testing"
)

const (
	SortAsc = 1
	SortDsc = 0
)

var (
	client     *mongo.Client
	updateOpt  = options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	replaceOpt = options.FindOneAndReplace().SetUpsert(true).SetReturnDocument(options.After)
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
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(config.Runtime.MongoDBUrl))
	if err != nil {
		log.Fatalf("unable to create mongodb client due to %v", err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatalf("unable to connect to mongodb at %v due to %v", config.Runtime.MongoDBUrl, err)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatalf("unable to ping mongodb at %v due to %v", config.Runtime.MongoDBUrl, err)
	}

	coll := client.Database(config.Runtime.MongoDBName).Collection("tickers")
	repos.TickerRepo = MongoTickerRepo{coll}
	repos.TickerRepo.Init()

	coll = client.Database(config.Runtime.MongoDBName).Collection("patterns")
	repos.PatternRepo = MongoPatternRepo{coll}
	repos.PatternRepo.Init()

	coll = client.Database(config.Runtime.MongoDBName).Collection("periods")
	repos.PeriodRepo = MongoPeriodRepo{c: coll}
	repos.PeriodRepo.Init()

	coll = client.Database(config.Runtime.MongoDBName).Collection("series")
	repos.SeriesRepo = MongoSeriesRepo{coll}
	repos.SeriesRepo.Init()

	repos.GraphController = &GraphController{repos.PeriodRepo, repos.PatternRepo}
}

func (repos *Repos) DropAll(t *testing.T) {

	if t == nil {
		log.Errorf("only able to delete databases when in test mode")
		return
	}
	db := client.Database(repos.config.Runtime.MongoDBName)
	err := db.Drop(context.TODO())
	if err != nil {
		log.Errorf("unable to drop database %v due to %v", repos.config.Runtime.MongoDBName, err)
	}
}

func createCollection(c *mongo.Collection, doc interface{}) (bool, error) {

	log.Infof("Checking if collection '%v' needs to be created", c.Name())
	created := false
	count, err := c.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return created, err
	}

	if count <= 0 {
		insertResult, err := c.InsertOne(context.TODO(), doc)
		if err != nil {
			return created, err
		}
		delResult, err := c.DeleteOne(context.TODO(), bson.D{{"_id", insertResult.InsertedID}})
		if err != nil {
			return created, err
		}
		if delResult.DeletedCount == 0 {
			return created, fmt.Errorf("problem creating collection '%v' with priming read", c.Name())
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
	InsertMany(data []*model.Pattern) (*mongo.InsertManyResult, error)
	DeleteByLength(length int) error
	DropAndCreate() error
	FindOneAndReplace(pattern *model.Pattern) *model.Pattern
	FindAndReplace(pattern *model.Pattern) *model.Pattern
	FindBySymbol(symbol string) ([]*model.Pattern, error)
	FindOneBySymbolAndValue(symbol, value string) (*model.Pattern, error)
	FindHighestUpProbability(density model.PatternDensity) (*model.Pattern, error)
	FindHighestDownProbability(density model.PatternDensity) (*model.Pattern, error)
	FindHighestNoChangeProbability(density model.PatternDensity) (*model.Pattern, error)
	FindLowestUpProbability(density model.PatternDensity) (*model.Pattern, error)
	FindLowestDownProbability(density model.PatternDensity) (*model.Pattern, error)
	FindLowestNoChangeProbability(density model.PatternDensity) (*model.Pattern, error)
}

type PeriodRepo interface {
	Init()
	InsertMany(data []*model.Period) (*mongo.InsertManyResult, error)
	DropAndCreate() error
	FindOneAndReplace(data *model.Period) *model.Period
	FindAndReplace(data *model.Period) *model.Period
	FindOneAndUpdateDailyResult(data *model.Period) (*model.Period, error)
	FindBySymbol(symbol string, sort int) (model.PeriodSlice, error)
	FindOneBySymbolAndValue(symbol, value string) (*model.Period, error)
}

type TickerRepo interface {
	Init()
	CountAll() (int64, error)
	InsertOne(ticker *model.Ticker) error
	InsertMany(data []*model.Ticker) error
	DropAndCreate() error
	FindOneAndReplace(ticker *model.Ticker) *model.Ticker
	FindAndReplace(ticker *model.Ticker) *model.Ticker
	FindOne(symbol string) (*model.Ticker, error)
	FindOneCompanyName(symbol string) (string, error)
	FindOneAndUpdateCompanyName(symbol, company string) *model.Ticker
	FindSymbols() []string
	FindSymbolsAndCompany() *report.TickerSymbolCompanySlice
}

type SeriesRepo interface {
	Init()
	FindBySymbol(symbol string) ([]model.Series, error)
	InsertOne(data *model.Series) error
	DeleteOne(data *model.Series) error
	DeleteByLength(length int) error
	DropAndCreate() error
}
