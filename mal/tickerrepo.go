package mal

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"market-patterns/model"
)

type TickerRepo struct {
	c          *mongo.Collection
	updateOpt  *options.FindOneAndUpdateOptions
	replaceOpt *options.FindOneAndReplaceOptions
}

var (
	TRUE  = true
	FALSE = false
)

func (repo *TickerRepo) Init() {

	repo.updateOpt = options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
	repo.replaceOpt = options.FindOneAndReplace().SetUpsert(true).SetReturnDocument(options.After)

	idxModel := mongo.IndexModel{}
	idxModel.Keys = bsonx.Doc{{Key: "symbol", Value: bsonx.Int32(1)}}

	name := "idx_symbol"
	idxModel.Options = &options.IndexOptions{Background: &TRUE, Name: &name, Unique: &TRUE}

	tmp, err := repo.c.Indexes().CreateOne(context.TODO(), idxModel)
	if err != nil {
		log.Errorf("problem creating %v due to %v", tmp, err)
	}
}

func (repo *TickerRepo) FindOneAndReplace(ticker *model.Ticker) *model.Ticker {

	filter := bson.D{{"symbol", ticker.Symbol}}

	var update model.Ticker
	err := repo.c.FindOneAndReplace(context.TODO(), filter, ticker, repo.replaceOpt).Decode(&update)
	if err != nil {
		log.Warnf("problem replacing ticker due to %v", err)
	}

	return &update
}

func (repo *TickerRepo) FindOne(symbol string) *model.Ticker {
	filter := bson.D{{"symbol", symbol}}

	var ticker model.Ticker
	err := repo.c.FindOne(context.TODO(), filter).Decode(&ticker)
	if err != nil {
		log.Warnf("unable to find ticker with symbol %v due to %v", symbol, err)
		return nil
	}

	return &ticker
}

func (repo *TickerRepo) FindSymbols() *[]string {

	var symbols []string

	ary, err := repo.c.Distinct(context.TODO(), "symbol", bson.D{})
	if err != nil {
		log.Warnf("unable to load ticker names due to %v", err)
		return &symbols
	}

	if len(ary) > 0 {
		symbols = make([]string, len(ary))
		for i, v := range ary {
			symbols[i] = fmt.Sprint(v)
		}
	}

	return &symbols
}
