package mal

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"market-patterns/model"
)

const (
	idxPeriodSymbolDate = "idxSymbolDate"
)

type MongoPeriodRepo struct {
	c *mongo.Collection
}

var (
	sortSymbolAsc = bson.D{{"symbol", 1}}
	sortSymbolDsc = bson.D{{"symbol", 0}}
)

func (repo MongoPeriodRepo) Init() {

	created, err := createCollection(repo.c, model.Period{})
	if err != nil {
		log.WithError(err).Fatal("Unable to continue initializing PeriodRepo")
	}

	if created {
		idxModel := mongo.IndexModel{}
		idxModel.Keys = bsonx.Doc{{Key: "symbol", Value: bsonx.Int32(1)}, {Key: "date", Value: bsonx.Int32(1)}}
		idxModel.Options = &options.IndexOptions{}
		idxModel.Options.SetUnique(true)
		idxModel.Options.SetName(idxPeriodSymbolDate)

		tmp, err := repo.c.Indexes().CreateOne(context.TODO(), idxModel)
		if err != nil {
			log.WithError(err).Errorf("problem creating '%v' index", tmp)
		} else {
			log.Infof("Created index '%v'", tmp)
		}
	}
}

// *********************************************************
//   Insert functions
// *********************************************************

func (repo MongoPeriodRepo) InsertMany(data []*model.Period) (*mongo.InsertManyResult, error) {

	dataAry := make([]interface{}, len(data))
	for i, v := range data {
		dataAry[i] = v
	}
	result, err := repo.c.InsertMany(context.TODO(), dataAry)
	if err != nil {
		return result, errors.Wrap(err, "problem inserting many periods")
	}
	return result, nil
}

// *********************************************************
//   Delete functions
// *********************************************************

func (repo MongoPeriodRepo) DropAndCreate() error {
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

func (repo MongoPeriodRepo) FindOneAndReplace(data *model.Period) *model.Period {

	filter := bson.D{{"symbol", data.Symbol}, {"date", data.Date}}
	var update model.Period
	err := repo.c.FindOneAndReplace(context.TODO(), filter, data, replaceOpt).Decode(&update)
	if err != nil {
		log.Warnf("problem replacing pattern due to %v", err)
	}
	return &update
}

func (repo MongoPeriodRepo) FindAndReplace(data *model.Period) *model.Period {

	filter := bson.D{{"symbol", data.Symbol}, {"date", data.Date}}
	var update model.Period
	err := repo.c.FindOneAndReplace(context.TODO(), filter, data, replaceOpt).Decode(&update)
	if err != nil {
		log.Warnf("problem replacing pattern due to %v", err)
	}
	return &update
}

func (repo MongoPeriodRepo) FindOneAndUpdateDailyResult(data *model.Period) (*model.Period, error) {

	filter := bson.D{{"symbol", data.Symbol}, {"date", data.Date}}
	update := bson.D{{"$set", bson.D{{"dailyResult", data.DailyResult}}}}
	var doc model.Period
	err := repo.c.FindOneAndUpdate(context.TODO(), filter, data, updateOpt).Decode(&update)
	if err != nil {
		return &doc, errors.Wrap(err, "problem updating period daily result")
	}
	return &doc, nil
}

func (repo MongoPeriodRepo) FindBySymbol(symbol string, sort int) (model.PeriodSlice, error) {

	opts := &options.FindOptions{}
	if sort == 1 {
		opts.Sort = sortSymbolAsc
	} else {
		opts.Sort = sortSymbolDsc
	}

	filter := bson.D{{"symbol", symbol}}
	var findData model.PeriodSlice
	cur, err := repo.c.Find(context.TODO(), filter, opts)
	if err != nil {
		return findData, errors.Wrap(err, "unable to find by symbol")
	}

	var results error

	for cur.Next(context.TODO()) {
		var doc model.Period
		err = cur.Decode(&doc)
		if err != nil {
			results = multierror.Append(results, err)
			continue
		}
		findData = append(findData, &doc)
	}
	return findData, results
}

func (repo MongoPeriodRepo) FindOneBySymbolAndValue(symbol, value string) (*model.Period, error) {

	filter := bson.D{{"symbol", symbol}, {"value", value}}

	var findData model.Period
	err := repo.c.FindOne(context.TODO(), filter).Decode(&findData)
	if err != nil {
		return &findData, errors.Wrapf(err, "unable to find pattern by symbol '%v' and value '%v", symbol, value)
	}
	return &findData, nil
}
