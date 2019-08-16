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
	idxSeriesSymbolLength = "idxSymbolLength"
)

type MongoSeriesRepo struct {
	c *mongo.Collection
}

func NewMongoSeriesRepo(c *mongo.Collection) *MongoSeriesRepo {
	return &MongoSeriesRepo{c}
}

func (repo MongoSeriesRepo) Init() {

	created, err := CreateCollection(repo.c, model.Series{})
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

func (repo MongoSeriesRepo) FindBySymbol(symbol string) ([]*model.Series, error) {
	filter := bson.D{{"symbol", symbol}}
	var findData []*model.Series
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
		findData = append(findData, &doc)
	}
	return findData, results
}

// *********************************************************
//   Insert functions
// *********************************************************

func (repo MongoSeriesRepo) InsertOne(data *model.Series) error {

	_, err := repo.c.InsertOne(context.TODO(), data)
	if err != nil {
		return errors.Wrap(err, "problem inserting one series")
	}
	return nil
}

// *********************************************************
//   Delete functions
// *********************************************************

func (repo MongoSeriesRepo) DeleteOne(data *model.Series) error {

	filter := bson.D{{"name", data.Name}}
	_, err := repo.c.DeleteOne(context.TODO(), filter)
	if err != nil {
		return errors.Wrap(err, "problem inserting one series")
	}
	return nil
}

func (repo MongoSeriesRepo) DeleteByLength(length int) error {

	filter := bson.D{{"length", length}}
	r, err := repo.c.DeleteMany(context.TODO(), filter)
	if err != nil {
		return errors.Wrapf(err, "problem deleting series with length %v", length)
	}

	log.Infof("Deleted %v docs with length %v from series repo", r.DeletedCount, length)
	return nil
}

func (repo MongoSeriesRepo) DropAndCreate() error {
	err := repo.c.Drop(context.TODO())
	if err != nil {
		return err
	}

	repo.Init()
	return nil
}
