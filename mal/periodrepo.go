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

type PeriodRepo struct {
	c          *mongo.Collection
	updateOpt  *options.FindOneAndUpdateOptions
	replaceOpt *options.FindOneAndReplaceOptions
	sortAsc    *bson.D
}

func (repo *PeriodRepo) SortAsc() *bson.D {
	return repo.sortAsc
}

func (repo *PeriodRepo) Init() {

	repo.updateOpt = options.FindOneAndUpdate().SetUpsert(TRUE).SetReturnDocument(options.After)
	repo.replaceOpt = options.FindOneAndReplace().SetUpsert(TRUE).SetReturnDocument(options.After)
	idxModel := mongo.IndexModel{}
	idxModel.Keys = bsonx.Doc{{Key: "symbol", Value: bsonx.Int32(1)}, {Key: "date", Value: bsonx.Int32(1)}}
	name := "idx_symbol_date"
	idxModel.Options = &options.IndexOptions{Background: &TRUE, Name: &name, Unique: &TRUE}
	tmp, err := repo.c.Indexes().CreateOne(context.TODO(), idxModel)
	if err != nil {
		log.Errorf("problem creating %v due to %v", tmp, err)
	}

	repo.sortAsc = &bson.D{{"symbol", 1}}
}

func (repo *PeriodRepo) InsertMany(data []*model.Period) error {

	dataAry := make([]interface{}, len(data))
	for i, v := range data {
		dataAry[i] = v
	}

	_, err := repo.c.InsertMany(context.TODO(), dataAry)
	if err != nil {
		return errors.Wrap(err, "problem inserting many periods")
	}
	return nil
}

func (repo *PeriodRepo) DeleteAll() error {
	return repo.c.Drop(context.TODO())
}

func (repo *PeriodRepo) FindOneAndReplace(data *model.Period) *model.Period {

	filter := bson.D{{"symbol", data.Symbol}, {"date", data.Date}}
	var update model.Period
	err := repo.c.FindOneAndReplace(context.TODO(), filter, data, repo.replaceOpt).Decode(&update)
	if err != nil {
		log.Warnf("problem replacing pattern due to %v", err)
	}
	return &update
}

func (repo *PeriodRepo) FindAndReplace(data *model.Period) *model.Period {

	filter := bson.D{{"symbol", data.Symbol}, {"date", data.Date}}
	var update model.Period
	err := repo.c.FindOneAndReplace(context.TODO(), filter, data, repo.replaceOpt).Decode(&update)
	if err != nil {
		log.Warnf("problem replacing pattern due to %v", err)
	}
	return &update
}

func (repo *PeriodRepo) FindOneAndUpdateDailyResult(data *model.Period) (*model.Period, error) {

	filter := bson.D{{"symbol", data.Symbol}, {"date", data.Date}}
	update := bson.D{{"$set", bson.D{{"dailyResult", data.DailyResult}}}}
	var doc model.Period
	err := repo.c.FindOneAndUpdate(context.TODO(), filter, data, repo.updateOpt).Decode(&update)
	if err != nil {
		return &doc, errors.Wrap(err, "problem updating period daily result")
	}
	return &doc, nil
}

func (repo *PeriodRepo) FindBySymbol(symbol string, sort *bson.D) (model.PeriodSlice, error) {

	opts := &options.FindOptions{}
	if sort != nil {
		opts.Sort = sort
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

func (repo *PeriodRepo) FindOneBySymbolAndValue(symbol, value string) (*model.Period, error) {

	filter := bson.D{{"symbol", symbol}, {"value", value}}

	var findData model.Period
	err := repo.c.FindOne(context.TODO(), filter).Decode(&findData)
	if err != nil {
		return &findData, errors.Wrapf(err, "unable to find pattern by symbol '%v' and value '%v", symbol, value)
	}
	return &findData, nil
}
