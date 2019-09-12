package mal

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go-market-patterns/model/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	idxPatternSymbolValueLength = "idxSymbolValueLength"
)

type MongoPatternRepo struct {
	c *mongo.Collection
}

func NewMongoPatternRepo(c *mongo.Collection) *MongoPatternRepo {
	return &MongoPatternRepo{c}
}

func (repo *MongoPatternRepo) Init() {

	created, err := CreateCollection(repo.c, core.Pattern{})
	if err != nil {
		log.WithError(err).Fatal("Unable to continue initializing MongoPatternRepo")
	}

	if created {
		idxModel := mongo.IndexModel{}
		idxModel.Keys = bsonx.Doc{{Key: "symbol", Value: bsonx.Int32(1)},
			{Key: "value", Value: bsonx.Int32(1)},
			{Key: "length", Value: bsonx.Int32(1)}}
		idxModel.Options = &options.IndexOptions{}
		idxModel.Options.SetUnique(true)
		idxModel.Options.SetName(idxPatternSymbolValueLength)

		tmp, err := repo.c.Indexes().CreateOne(context.TODO(), idxModel)
		if err != nil {
			log.WithError(err).Errorf("problem creating '%v' index", tmp)
		} else {
			log.Infof("Created index '%v'", tmp)
		}
	}
}

// *********************************************************
// Insert functions
// *********************************************************

func (repo *MongoPatternRepo) InsertMany(data []*core.Pattern) (int, error) {

	dataAry := make([]interface{}, len(data))
	for i, v := range data {
		dataAry[i] = v
	}
	results, err := repo.c.InsertMany(context.TODO(), dataAry)
	if err != nil {
		return len(results.InsertedIDs), errors.Wrap(err, "problem inserting many patterns")
	}
	return len(results.InsertedIDs), nil
}

// *********************************************************
// Delete functions
// *********************************************************

func (repo *MongoPatternRepo) DeleteByLength(length int) error {

	filter := bson.D{{"length", length}}
	r, err := repo.c.DeleteMany(context.TODO(), filter)
	if err != nil {
		return errors.Wrapf(err, "problem deleting patterns with series length %v", length)
	}

	log.Infof("Deleted %v docs with series length %v from patterns repo", r.DeletedCount, length)
	return nil
}

func (repo *MongoPatternRepo) DropAndCreate() error {
	err := repo.c.Drop(context.TODO())
	if err != nil {
		return err
	}

	repo.Init()
	return nil
}

// *********************************************************
// Find functions
// *********************************************************

func (repo *MongoPatternRepo) FindOneAndReplace(pattern *core.Pattern) *core.Pattern {

	filter := bson.D{{"symbol", pattern.Symbol}, {"value", pattern.Value}}
	var update core.Pattern
	err := repo.c.FindOneAndReplace(context.TODO(), filter, pattern, ReplaceOpt).Decode(&update)
	if err != nil {
		log.Warnf("problem replacing pattern due to %v", err)
	}
	return &update
}

func (repo *MongoPatternRepo) FindAndReplace(pattern *core.Pattern) *core.Pattern {

	filter := bson.D{{"symbol", pattern.Symbol}, {"value", pattern.Value}}
	var update core.Pattern
	err := repo.c.FindOneAndReplace(context.TODO(), filter, pattern, ReplaceOpt).Decode(&update)
	if err != nil {
		log.Warnf("problem replacing pattern due to %v", err)
	}
	return &update
}

func (repo *MongoPatternRepo) FindBySymbol(symbol string) ([]*core.Pattern, error) {

	filter := bson.D{{"symbol", symbol}}
	var findData []*core.Pattern
	cur, err := repo.c.Find(context.TODO(), filter)
	if err != nil {
		return findData, errors.Wrap(err, "unable to find by symbol")
	}
	defer func(c *mongo.Cursor) {
		err := c.Close(context.TODO())
		if err != nil {
			log.Errorf("problem closing cursor due to %v", err)
		}
	}(cur)

	var results error

	for cur.Next(context.TODO()) {
		var doc core.Pattern
		err = cur.Decode(&doc)
		if err != nil {
			results = multierror.Append(results, err)
			continue
		}
		findData = append(findData, &doc)
	}
	return findData, results
}

func (repo *MongoPatternRepo) FindBySymbolAndLength(symbol string, length int) ([]*core.Pattern, error) {

	filter := bson.D{{"symbol", symbol}, {"length", length}}
	var findData []*core.Pattern
	cur, err := repo.c.Find(context.TODO(), filter)
	if err != nil {
		return findData, errors.Wrap(err, "unable to find by symbol")
	}
	defer func(c *mongo.Cursor) {
		err := c.Close(context.TODO())
		if err != nil {
			log.Errorf("problem closing cursor due to %v", err)
		}
	}(cur)

	var results error

	for cur.Next(context.TODO()) {
		var doc core.Pattern
		err = cur.Decode(&doc)
		if err != nil {
			results = multierror.Append(results, err)
			continue
		}
		findData = append(findData, &doc)
	}
	return findData, results
}

func (repo *MongoPatternRepo) FindOneBySymbolAndValueAndLength(symbol, value string, length int) (*core.Pattern, error) {

	filter := bson.D{{"symbol", symbol}, {"value", value}}

	var pattern core.Pattern
	err := repo.c.FindOne(context.TODO(), filter).Decode(&pattern)
	if err != nil {
		return &pattern, errors.Wrapf(err, "unable to find pattern by symbol '%v' and value '%v", symbol, value)
	}
	return &pattern, nil
}

var (
	patternAggregateUpMaxLowDensity = mongo.Pipeline{
		{{"$group", bson.D{{"_id", "$$ROOT"},
			{"max", bson.D{{"$max",
				bson.D{{"$divide", bson.A{"$upcount", "$totalcount"}}}}}}}}},
		{{"$sort", bson.D{{"max", -1}}}},
		{{"$limit", 1}},
		{{"$replaceRoot",
			bson.D{{"newRoot", "$_id"}}}},
	}

	patternAggregateUpMaxMediumDensity = mongo.Pipeline{
		{{"$match", bson.D{{"totalcount", bson.D{{"$gte", 500}}}}}},
		{{"$group", bson.D{{"_id", "$$ROOT"},
			{"max", bson.D{{"$max",
				bson.D{{"$divide", bson.A{"$upcount", "$totalcount"}}}}}}}}},
		{{"$sort", bson.D{{"max", -1}}}},
		{{"$limit", 1}},
		{{"$replaceRoot",
			bson.D{{"newRoot", "$_id"}}}},
	}

	patternAggregateUpMaxHighDensity = mongo.Pipeline{
		{{"$match", bson.D{{"totalcount", bson.D{{"$gte", 1000}}}}}},
		{{"$group", bson.D{{"_id", "$$ROOT"},
			{"max", bson.D{{"$max",
				bson.D{{"$divide", bson.A{"$upcount", "$totalcount"}}}}}}}}},
		{{"$sort", bson.D{{"max", -1}}}},
		{{"$limit", 1}},
		{{"$replaceRoot",
			bson.D{{"newRoot", "$_id"}}}},
	}

	patternAggregateDownMaxLowDensity = mongo.Pipeline{
		{{"$group", bson.D{{"_id", "$$ROOT"},
			{"max", bson.D{{"$max",
				bson.D{{"$divide", bson.A{"$downcount", "$totalcount"}}}}}}}}},
		{{"$sort", bson.D{{"max", -1}}}},
		{{"$limit", 1}},
		{{"$replaceRoot",
			bson.D{{"newRoot", "$_id"}}}},
	}

	patternAggregateDownMaxMediumDensity = mongo.Pipeline{
		{{"$match", bson.D{{"totalcount", bson.D{{"$gte", 500}}}}}},
		{{"$group", bson.D{{"_id", "$$ROOT"},
			{"max", bson.D{{"$max",
				bson.D{{"$divide", bson.A{"$downcount", "$totalcount"}}}}}}}}},
		{{"$sort", bson.D{{"max", -1}}}},
		{{"$limit", 1}},
		{{"$replaceRoot",
			bson.D{{"newRoot", "$_id"}}}},
	}

	patternAggregateDownMaxHighDensity = mongo.Pipeline{
		{{"$match", bson.D{{"totalcount", bson.D{{"$gte", 1000}}}}}},
		{{"$group", bson.D{{"_id", "$$ROOT"},
			{"max", bson.D{{"$max",
				bson.D{{"$divide", bson.A{"$downcount", "$totalcount"}}}}}}}}},
		{{"$sort", bson.D{{"max", -1}}}},
		{{"$limit", 1}},
		{{"$replaceRoot",
			bson.D{{"newRoot", "$_id"}}}},
	}

	patternAggregateNoChangeMaxLowDensity = mongo.Pipeline{
		{{"$group", bson.D{{"_id", "$$ROOT"},
			{"max", bson.D{{"$max",
				bson.D{{"$divide", bson.A{"$nochangecount", "$totalcount"}}}}}}}}},
		{{"$sort", bson.D{{"max", -1}}}},
		{{"$limit", 1}},
		{{"$replaceRoot",
			bson.D{{"newRoot", "$_id"}}}},
	}

	patternAggregateNoChangeMaxMediumDensity = mongo.Pipeline{
		{{"$match", bson.D{{"totalcount", bson.D{{"$gte", 500}}}}}},
		{{"$group", bson.D{{"_id", "$$ROOT"},
			{"max", bson.D{{"$max",
				bson.D{{"$divide", bson.A{"$nochangecount", "$totalcount"}}}}}}}}},
		{{"$sort", bson.D{{"max", -1}}}},
		{{"$limit", 1}},
		{{"$replaceRoot",
			bson.D{{"newRoot", "$_id"}}}},
	}

	patternAggregateNoChangeMaxHighDensity = mongo.Pipeline{
		{{"$match", bson.D{{"totalcount", bson.D{{"$gte", 1000}}}}}},
		{{"$group", bson.D{{"_id", "$$ROOT"},
			{"max", bson.D{{"$max",
				bson.D{{"$divide", bson.A{"$nochangecount", "$totalcount"}}}}}}}}},
		{{"$sort", bson.D{{"max", -1}}}},
		{{"$limit", 1}},
		{{"$replaceRoot",
			bson.D{{"newRoot", "$_id"}}}},
	}

	patternAggregateUpMinLowDensity = mongo.Pipeline{
		{{"$group", bson.D{{"_id", "$$ROOT"},
			{"max", bson.D{{"$max",
				bson.D{{"$divide", bson.A{"$upcount", "$totalcount"}}}}}}}}},
		{{"$sort", bson.D{{"max", -1}}}},
		{{"$limit", 1}},
		{{"$replaceRoot",
			bson.D{{"newRoot", "$_id"}}}},
	}

	patternAggregateUpMinMediumDensity = mongo.Pipeline{
		{{"$match", bson.D{{"totalcount", bson.D{{"$gte", 500}}}}}},
		{{"$group", bson.D{{"_id", "$$ROOT"},
			{"max", bson.D{{"$max",
				bson.D{{"$divide", bson.A{"$upcount", "$totalcount"}}}}}}}}},
		{{"$sort", bson.D{{"max", -1}}}},
		{{"$limit", 1}},
		{{"$replaceRoot",
			bson.D{{"newRoot", "$_id"}}}},
	}

	patternAggregateUpMinHighDensity = mongo.Pipeline{
		{{"$match", bson.D{{"totalcount", bson.D{{"$gte", 1000}}}}}},
		{{"$group", bson.D{{"_id", "$$ROOT"},
			{"max", bson.D{{"$max",
				bson.D{{"$divide", bson.A{"$upcount", "$totalcount"}}}}}}}}},
		{{"$sort", bson.D{{"max", -1}}}},
		{{"$limit", 1}},
		{{"$replaceRoot",
			bson.D{{"newRoot", "$_id"}}}},
	}

	patternAggregateDownMinLowDensity = mongo.Pipeline{
		{{"$group", bson.D{{"_id", "$$ROOT"},
			{"max", bson.D{{"$max",
				bson.D{{"$divide", bson.A{"$downcount", "$totalcount"}}}}}}}}},
		{{"$sort", bson.D{{"max", -1}}}},
		{{"$limit", 1}},
		{{"$replaceRoot",
			bson.D{{"newRoot", "$_id"}}}},
	}

	patternAggregateDownMinMediumDensity = mongo.Pipeline{
		{{"$match", bson.D{{"totalcount", bson.D{{"$gte", 500}}}}}},
		{{"$group", bson.D{{"_id", "$$ROOT"},
			{"max", bson.D{{"$max",
				bson.D{{"$divide", bson.A{"$downcount", "$totalcount"}}}}}}}}},
		{{"$sort", bson.D{{"max", -1}}}},
		{{"$limit", 1}},
		{{"$replaceRoot",
			bson.D{{"newRoot", "$_id"}}}},
	}

	patternAggregateDownMinHighDensity = mongo.Pipeline{
		{{"$match", bson.D{{"totalcount", bson.D{{"$gte", 1000}}}}}},
		{{"$group", bson.D{{"_id", "$$ROOT"},
			{"max", bson.D{{"$max",
				bson.D{{"$divide", bson.A{"$downcount", "$totalcount"}}}}}}}}},
		{{"$sort", bson.D{{"max", -1}}}},
		{{"$limit", 1}},
		{{"$replaceRoot",
			bson.D{{"newRoot", "$_id"}}}},
	}

	patternAggregateNoChangeMinLowDensity = mongo.Pipeline{
		{{"$group", bson.D{{"_id", "$$ROOT"},
			{"max", bson.D{{"$max",
				bson.D{{"$divide", bson.A{"$nochangecount", "$totalcount"}}}}}}}}},
		{{"$sort", bson.D{{"max", -1}}}},
		{{"$limit", 1}},
		{{"$replaceRoot",
			bson.D{{"newRoot", "$_id"}}}},
	}

	patternAggregateNoChangeMinMediumDensity = mongo.Pipeline{
		{{"$match", bson.D{{"totalcount", bson.D{{"$gte", 500}}}}}},
		{{"$group", bson.D{{"_id", "$$ROOT"},
			{"max", bson.D{{"$max",
				bson.D{{"$divide", bson.A{"$nochangecount", "$totalcount"}}}}}}}}},
		{{"$sort", bson.D{{"max", -1}}}},
		{{"$limit", 1}},
		{{"$replaceRoot",
			bson.D{{"newRoot", "$_id"}}}},
	}

	patternAggregateNoChangeMinHighDensity = mongo.Pipeline{
		{{"$match", bson.D{{"totalcount", bson.D{{"$gte", 1000}}}}}},
		{{"$group", bson.D{{"_id", "$$ROOT"},
			{"max", bson.D{{"$max",
				bson.D{{"$divide", bson.A{"$nochangecount", "$totalcount"}}}}}}}}},
		{{"$sort", bson.D{{"max", -1}}}},
		{{"$limit", 1}},
		{{"$replaceRoot",
			bson.D{{"newRoot", "$_id"}}}},
	}
)

func (repo *MongoPatternRepo) FindHighestUpProbability(density core.PatternDensity) (*core.Pattern, error) {

	var pipeline mongo.Pipeline
	switch density {
	case core.PatternDensityLow:
		pipeline = patternAggregateUpMaxLowDensity
	case core.PatternDensityMedium:
		pipeline = patternAggregateUpMaxMediumDensity
	case core.PatternDensityHigh:
		pipeline = patternAggregateUpMaxHighDensity
	}

	cur, err := repo.c.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer func(c *mongo.Cursor) {
		err := c.Close(context.TODO())
		if err != nil {
			log.Errorf("problem closing cursor due to %v", err)
		}
	}(cur)

	var pattern core.Pattern
	// Only get the first one
	if cur.Next(context.TODO()) {
		err := cur.Decode(&pattern)
		if err != nil {
			return &pattern, err
		}
	}

	return &pattern, nil
}

func (repo *MongoPatternRepo) FindHighestDownProbability(density core.PatternDensity) (*core.Pattern, error) {

	var pipeline mongo.Pipeline
	switch density {
	case core.PatternDensityLow:
		pipeline = patternAggregateDownMaxLowDensity
	case core.PatternDensityMedium:
		pipeline = patternAggregateDownMaxMediumDensity
	case core.PatternDensityHigh:
		pipeline = patternAggregateDownMaxHighDensity
	}

	cur, err := repo.c.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer func(c *mongo.Cursor) {
		err := c.Close(context.TODO())
		if err != nil {
			log.Errorf("problem closing cursor due to %v", err)
		}
	}(cur)

	var pattern core.Pattern
	// Only get the first one
	if cur.Next(context.TODO()) {
		err := cur.Decode(&pattern)
		if err != nil {
			return &pattern, err
		}
	}

	return &pattern, nil
}

func (repo *MongoPatternRepo) FindHighestNoChangeProbability(density core.PatternDensity) (*core.Pattern, error) {

	var pipeline mongo.Pipeline
	switch density {
	case core.PatternDensityLow:
		pipeline = patternAggregateNoChangeMaxLowDensity
	case core.PatternDensityMedium:
		pipeline = patternAggregateNoChangeMaxMediumDensity
	case core.PatternDensityHigh:
		pipeline = patternAggregateNoChangeMaxHighDensity
	}

	cur, err := repo.c.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer func(c *mongo.Cursor) {
		err := c.Close(context.TODO())
		if err != nil {
			log.Errorf("problem closing cursor due to %v", err)
		}
	}(cur)

	var pattern core.Pattern
	// Only get the first one
	if cur.Next(context.TODO()) {
		err := cur.Decode(&pattern)
		if err != nil {
			return &pattern, err
		}
	}

	return &pattern, nil
}

func (repo *MongoPatternRepo) FindLowestUpProbability(density core.PatternDensity) (*core.Pattern, error) {

	var pipeline mongo.Pipeline
	switch density {
	case core.PatternDensityLow:
		pipeline = patternAggregateUpMinLowDensity
	case core.PatternDensityMedium:
		pipeline = patternAggregateUpMinMediumDensity
	case core.PatternDensityHigh:
		pipeline = patternAggregateUpMinHighDensity
	}

	cur, err := repo.c.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer func(c *mongo.Cursor) {
		err := c.Close(context.TODO())
		if err != nil {
			log.Errorf("problem closing cursor due to %v", err)
		}
	}(cur)

	var pattern core.Pattern
	// Only get the first one
	if cur.Next(context.TODO()) {
		err := cur.Decode(&pattern)
		if err != nil {
			return &pattern, err
		}
	}

	return &pattern, nil
}

func (repo *MongoPatternRepo) FindLowestDownProbability(density core.PatternDensity) (*core.Pattern, error) {

	var pipeline mongo.Pipeline
	switch density {
	case core.PatternDensityLow:
		pipeline = patternAggregateDownMinLowDensity
	case core.PatternDensityMedium:
		pipeline = patternAggregateDownMinMediumDensity
	case core.PatternDensityHigh:
		pipeline = patternAggregateDownMinHighDensity
	}

	cur, err := repo.c.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer func(c *mongo.Cursor) {
		err := c.Close(context.TODO())
		if err != nil {
			log.Errorf("problem closing cursor due to %v", err)
		}
	}(cur)

	var pattern core.Pattern
	// Only get the first one
	if cur.Next(context.TODO()) {
		err := cur.Decode(&pattern)
		if err != nil {
			return &pattern, err
		}
	}

	return &pattern, nil
}

func (repo *MongoPatternRepo) FindLowestNoChangeProbability(density core.PatternDensity) (*core.Pattern, error) {

	var pipeline mongo.Pipeline
	switch density {
	case core.PatternDensityLow:
		pipeline = patternAggregateNoChangeMinLowDensity
	case core.PatternDensityMedium:
		pipeline = patternAggregateNoChangeMinMediumDensity
	case core.PatternDensityHigh:
		pipeline = patternAggregateNoChangeMinHighDensity
	}

	cur, err := repo.c.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer func(c *mongo.Cursor) {
		err := c.Close(context.TODO())
		if err != nil {
			log.Errorf("problem closing cursor due to %v", err)
		}
	}(cur)

	var pattern core.Pattern
	// Only get the first one
	if cur.Next(context.TODO()) {
		err := cur.Decode(&pattern)
		if err != nil {
			return &pattern, err
		}
	}

	return &pattern, nil
}
