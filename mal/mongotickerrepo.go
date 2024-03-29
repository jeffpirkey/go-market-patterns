package mal

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go-market-patterns/model"
	"go-market-patterns/model/report"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"sort"
)

const (
	idxTickerSymbol = "idxSymbol"
)

type MongoTickerRepo struct {
	c *mongo.Collection
}

func NewMongoTickerRepo(c *mongo.Collection) *MongoTickerRepo {
	return &MongoTickerRepo{c}
}

func (repo *MongoTickerRepo) Init() {
	created, err := CreateCollection(repo.c, model.Ticker{})
	if err != nil {
		log.WithError(err).Fatal("Unable to continue initializing MongoTickerRepo")
	}

	if created {
		idxModel := mongo.IndexModel{}
		idxModel.Keys = bsonx.Doc{{Key: "symbol", Value: bsonx.Int32(1)}}
		idxModel.Options = &options.IndexOptions{}
		idxModel.Options.SetUnique(true)
		idxModel.Options.SetName(idxTickerSymbol)

		tmp, err := repo.c.Indexes().CreateOne(context.TODO(), idxModel)
		if err != nil {
			log.WithError(err).Errorf("problem creating '%v' index", tmp)
		} else {
			log.Infof("Created index '%v'", tmp)
		}
	}
}

func (repo *MongoTickerRepo) CountAll() (int64, error) {
	return repo.c.CountDocuments(context.TODO(), bson.D{})
}

// *********************************************************
//   Insert functions
// *********************************************************

func (repo *MongoTickerRepo) InsertOne(ticker *model.Ticker) error {

	_, err := repo.c.InsertOne(context.TODO(), ticker)
	if err != nil {
		return errors.Wrap(err, "problem inserting ticker")
	}

	return nil
}

func (repo *MongoTickerRepo) InsertMany(data []*model.Ticker) error {

	dataAry := make([]interface{}, len(data))
	for i, v := range data {
		dataAry[i] = v
	}
	_, err := repo.c.InsertMany(context.TODO(), dataAry)
	if err != nil {
		return errors.Wrap(err, "problem inserting many tickers")
	}
	return nil
}

func (repo *MongoTickerRepo) DropAndCreate() error {
	err := repo.c.Drop(context.TODO())
	if err != nil {
		return err
	}

	repo.Init()
	return nil
}

// *********************************************************
//   Find functions
// *********************************************************

func (repo *MongoTickerRepo) FindOne(symbol string) (*model.Ticker, error) {

	filter := bson.D{{"symbol", symbol}}
	var ticker model.Ticker
	err := repo.c.FindOne(context.TODO(), filter).Decode(&ticker)
	return &ticker, err
}

func (repo *MongoTickerRepo) FindOneCompanyNameBySymbol(symbol string) (string, error) {

	filter := bson.D{{"symbol", symbol}}
	var ticker model.Ticker
	err := repo.c.FindOne(context.TODO(), filter).Decode(&ticker)
	return ticker.Company, err
}

func (repo *MongoTickerRepo) FindOneAndUpdateCompanyName(symbol, company string) *model.Ticker {
	filter := bson.D{{"symbol", symbol}}
	update := bson.D{{"$set", bson.D{{"company", company}}}}

	var result model.Ticker
	err := repo.c.FindOneAndUpdate(context.TODO(), filter, update, UpdateOpt).Decode(&result)
	if err != nil {
		log.Warnf("unable to update company of ticker with symbol %v due to %v", symbol, err)
		return nil
	}

	return &result
}

func (repo *MongoTickerRepo) FindSymbols() []string {

	var symbols []string

	ary, err := repo.c.Distinct(context.TODO(), "symbol", bson.D{})
	if err != nil {
		log.Warnf("unable to load ticker symbols due to %v", err)
		return symbols
	}

	if len(ary) > 0 {
		symbols = make([]string, len(ary))
		for i, v := range ary {
			symbols[i] = fmt.Sprint(v)
		}
	}

	return symbols
}

func (repo *MongoTickerRepo) FindSymbolsAndCompany() *report.TickerSymbolCompanySlice {

	var symbols report.TickerSymbolCompanySlice

	opts := options.Find()
	opts.Projection = bson.D{{"symbol", 1}, {"company", 1}, {"_id", 0}}
	cur, err := repo.c.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Warnf("unable to load ticker symbols and companies due to %v", err)
		return &symbols
	}

	for cur.Next(context.TODO()) {
		var doc *report.TickerSymbolCompany
		err := cur.Decode(&doc)
		if err != nil {
			log.Errorf("unable to unmarshal due to %v", err)
			continue
		}
		symbols = append(symbols, doc)
	}

	sort.Sort(symbols)

	return &symbols
}
