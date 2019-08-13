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
	"testing"
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
	TickerRepo      *TickerRepo
	PatternRepo     *PatternRepo
	PeriodRepo      *PeriodRepo
	SeriesRepo      *SeriesRepo
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
	repos.TickerRepo = &TickerRepo{coll}
	repos.TickerRepo.Init()

	coll = client.Database(config.Runtime.MongoDBName).Collection("patterns")
	repos.PatternRepo = &PatternRepo{coll}
	repos.PatternRepo.Init()

	coll = client.Database(config.Runtime.MongoDBName).Collection("periods")
	repos.PeriodRepo = &PeriodRepo{c: coll}
	repos.PeriodRepo.Init()

	coll = client.Database(config.Runtime.MongoDBName).Collection("series")
	repos.SeriesRepo = &SeriesRepo{coll}
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
