package mal

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"market-patterns/model"
	"market-patterns/model/report"
	"sort"
)

const (
	idxTickerSymbol = "idxSymbol"
)

type TickerRepo struct {
	c *mongo.Collection
}

func (repo *TickerRepo) Init() {

	created, err := createCollection(repo.c, model.Ticker{})
	if err != nil {
		log.WithError(err).Fatal("Unable to continue initializing TickerRepo")
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

func (repo *TickerRepo) CountAll() (int64, error) {
	return repo.c.CountDocuments(context.TODO(), bson.D{})
}

// *********************************************************
//   Insert functions
// *********************************************************

func (repo *TickerRepo) InsertOne(ticker *model.Ticker) error {

	_, err := repo.c.InsertOne(context.TODO(), ticker)
	if err != nil {
		return errors.Wrap(err, "problem inserting ticker")
	}

	return nil
}

func (repo *TickerRepo) InsertMany(data []*model.Ticker) error {

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

func (repo *TickerRepo) DropAndCreate() error {
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

func (repo *TickerRepo) FindOneAndReplace(ticker *model.Ticker) *model.Ticker {

	filter := bson.D{{"symbol", ticker.Symbol}}

	var update model.Ticker
	err := repo.c.FindOneAndReplace(context.TODO(), filter, ticker, replaceOpt).Decode(&update)
	if err != nil {
		log.Warnf("problem replacing ticker due to %v", err)
	}

	return &update
}

func (repo *TickerRepo) FindAndReplace(ticker *model.Ticker) *model.Ticker {

	filter := bson.D{{"symbol", ticker.Symbol}}

	var update model.Ticker
	err := repo.c.FindOneAndReplace(context.TODO(), filter, ticker, replaceOpt).Decode(&update)
	if err != nil {
		log.Warnf("problem replacing ticker due to %v", err)
	}

	return &update
}

func (repo *TickerRepo) FindOne(symbol string) (*model.Ticker, error) {

	filter := bson.D{{"symbol", symbol}}
	var ticker model.Ticker
	err := repo.c.FindOne(context.TODO(), filter).Decode(&ticker)
	return &ticker, err
}

func (repo *TickerRepo) FindOneCompanyName(symbol string) (string, error) {

	filter := bson.D{{"symbol", symbol}}
	var ticker model.Ticker
	err := repo.c.FindOne(context.TODO(), filter).Decode(&ticker)
	return ticker.Company, err
}

func (repo *TickerRepo) FindOneAndUpdateCompanyName(symbol, company string) *model.Ticker {
	filter := bson.D{{"symbol", symbol}}
	update := bson.D{{"$set", bson.D{{"company", company}}}}

	var result model.Ticker
	err := repo.c.FindOneAndUpdate(context.TODO(), filter, update, updateOpt).Decode(&result)
	if err != nil {
		log.Warnf("unable to update company of ticker with symbol %v due to %v", symbol, err)
		return nil
	}

	return &result
}

func (repo *TickerRepo) FindSymbols() []string {

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

func (repo *TickerRepo) FindSymbolsAndCompany() *report.TickerSymbolCompanySlice {

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
