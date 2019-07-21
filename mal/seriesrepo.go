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

type SeriesRepo struct {
	c          *mongo.Collection
	updateOpt  *options.FindOneAndUpdateOptions
	replaceOpt *options.FindOneAndReplaceOptions
}

func (repo *SeriesRepo) Init() {

	repo.updateOpt = options.FindOneAndUpdate().SetUpsert(TRUE).SetReturnDocument(options.After)
	repo.replaceOpt = options.FindOneAndReplace().SetUpsert(TRUE).SetReturnDocument(options.After)
	idxModel := mongo.IndexModel{}
	idxModel.Keys = bsonx.Doc{{Key: "symbol", Value: bsonx.Int32(1)}, {Key: "name", Value: bsonx.Int32(1)}}
	name := "idx_symbol_name"
	idxModel.Options = &options.IndexOptions{Background: &TRUE, Name: &name, Unique: &TRUE}
	tmp, err := repo.c.Indexes().CreateOne(context.TODO(), idxModel)
	if err != nil {
		log.Errorf("problem creating %v due to %v", tmp, err)
	}
}

func (repo *SeriesRepo) FindBySymbol(symbol string) ([]model.Series, error) {
	filter := bson.D{{"symbol", symbol}}
	var findData []model.Series
	cur, err := repo.c.Find(context.TODO(), filter)
	if err != nil {
		return findData, errors.Wrap(err, "unable to find by symbol")
	}

	var results error

	for cur.Next(context.TODO()) {
		var doc model.Series
		err = cur.Decode(&doc)
		if err != nil {
			results = multierror.Append(results, err)
			continue
		}
		findData = append(findData, doc)
	}
	return findData, results
}

func (repo *SeriesRepo) InsertOne(data *model.Series) error {

	_, err := repo.c.InsertOne(context.TODO(), data)
	if err != nil {
		return errors.Wrap(err, "problem inserting one series")
	}
	return nil
}

func (repo *SeriesRepo) DeleteOne(data *model.Series) error {

	filter := bson.D{{"name", data.Name}}
	_, err := repo.c.DeleteOne(context.TODO(), filter)
	if err != nil {
		return errors.Wrap(err, "problem inserting one series")
	}
	return nil
}

func (repo *SeriesRepo) DeleteAll() error {
	return repo.c.Drop(context.TODO())
}
