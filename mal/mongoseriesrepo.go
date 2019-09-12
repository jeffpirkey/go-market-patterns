package mal

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go-market-patterns/model/core"
	"go-market-patterns/model/report"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"sort"
)

const (
	idxSeriesSymbolLength = "idxSymbolLength"
)

type MongoSeriesRepo struct {
	c *mongo.Collection
}

func NewMongoSeriesRepo(c *mongo.Collection) *MongoSeriesRepo {
	return &MongoSeriesRepo{c}
}

func (repo *MongoSeriesRepo) Init() {

	created, err := CreateCollection(repo.c, core.Series{})
	if err != nil {
		log.WithError(err).Fatal("Unable to continue initializing MongoSeriesRepo")
	}

	if created {
		idxModel := mongo.IndexModel{}
		idxModel.Keys = bsonx.Doc{{Key: "symbol", Value: bsonx.Int32(1)},
			{Key: "length", Value: bsonx.Int32(1)}}
		idxModel.Options = &options.IndexOptions{}
		idxModel.Options.SetUnique(true)
		idxModel.Options.SetName(idxSeriesSymbolLength)

		tmp, err := repo.c.Indexes().CreateOne(context.TODO(), idxModel)
		if err != nil {
			log.WithError(err).Errorf("problem creating '%v' index", tmp)
		} else {
			log.Infof("Created index '%v'", tmp)
		}
	}
}

// *********************************************************
//   Find functions
// *********************************************************

func (repo *MongoSeriesRepo) FindBySymbol(symbol string) ([]*core.Series, error) {
	filter := bson.D{{"symbol", symbol}}
	var findData []*core.Series
	cur, err := repo.c.Find(context.TODO(), filter)
	if err != nil {
		return findData, errors.Wrap(err, "unable to find by symbol")
	}

	var results error

	for cur.Next(context.TODO()) {
		var doc core.Series
		err = cur.Decode(&doc)
		if err != nil {
			results = multierror.Append(results, err)
			continue
		}
		findData = append(findData, &doc)
	}
	return findData, results
}

func (repo *MongoSeriesRepo) FindOneBySymbolAndLength(symbol string, length int) (*core.Series, error) {
	filter := bson.D{{"symbol", symbol}, {"length", length}}
	var findData *core.Series
	result := repo.c.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return findData, errors.Wrap(result.Err(), "unable to find by symbol and length")
	}

	err := result.Decode(findData)
	return findData, err
}

func (repo *MongoSeriesRepo) FindNameLengthSliceBySymbol(symbol string) *report.SeriesNameLengthSlice {
	filter := bson.D{{"symbol", symbol}}
	var findData report.SeriesNameLengthSlice
	cur, err := repo.c.Find(context.TODO(), filter)
	if err != nil {
		log.Warnf("unable to load series names and lengths due to %v", err)
		return &findData
	}

	var results error
	for cur.Next(context.TODO()) {
		var doc report.SeriesNameLength
		err = cur.Decode(&doc)
		if err != nil {
			results = multierror.Append(results, err)
			continue
		}
		findData = append(findData, &doc)
	}
	if results != nil {
		log.Error(results)
	}

	sort.Sort(findData)

	return &findData
}

// *********************************************************
//   Insert functions
// *********************************************************

func (repo *MongoSeriesRepo) InsertOne(data *core.Series) error {

	_, err := repo.c.InsertOne(context.TODO(), data)
	if err != nil {
		return errors.Wrap(err, "problem inserting one series")
	}
	return nil
}

// *********************************************************
//   Delete functions
// *********************************************************

func (repo *MongoSeriesRepo) DeleteOne(data *core.Series) error {

	filter := bson.D{{"name", data.Name}}
	_, err := repo.c.DeleteOne(context.TODO(), filter)
	if err != nil {
		return errors.Wrap(err, "problem inserting one series")
	}
	return nil
}

func (repo *MongoSeriesRepo) DeleteByLength(length int) error {

	filter := bson.D{{"length", length}}
	r, err := repo.c.DeleteMany(context.TODO(), filter)
	if err != nil {
		return errors.Wrapf(err, "problem deleting series with length %v", length)
	}

	log.Infof("Deleted %v docs with length %v from series repo", r.DeletedCount, length)
	return nil
}

func (repo *MongoSeriesRepo) DropAndCreate() error {
	err := repo.c.Drop(context.TODO())
	if err != nil {
		return err
	}

	repo.Init()
	return nil
}
