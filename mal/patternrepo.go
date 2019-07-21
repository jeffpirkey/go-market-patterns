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

type PatternRepo struct {
	c          *mongo.Collection
	updateOpt  *options.FindOneAndUpdateOptions
	replaceOpt *options.FindOneAndReplaceOptions
}

func (repo *PatternRepo) Init() {

	repo.updateOpt = options.FindOneAndUpdate().SetUpsert(TRUE).SetReturnDocument(options.After)
	repo.replaceOpt = options.FindOneAndReplace().SetUpsert(TRUE).SetReturnDocument(options.After)
	idxModel := mongo.IndexModel{}
	idxModel.Keys = bsonx.Doc{{Key: "symbol", Value: bsonx.Int32(1)}, {Key: "value", Value: bsonx.Int32(1)}}
	name := "idx_symbol_date"
	idxModel.Options = &options.IndexOptions{Background: &TRUE, Name: &name, Unique: &TRUE}
	tmp, err := repo.c.Indexes().CreateOne(context.TODO(), idxModel)
	if err != nil {
		log.Errorf("problem creating %v due to %v", tmp, err)
	}
}

// *********************************************************
// Insert functions
// *********************************************************

func (repo *PatternRepo) InsertMany(data []*model.Pattern) error {

	dataAry := make([]interface{}, len(data))
	for i, v := range data {
		dataAry[i] = v
	}

	_, err := repo.c.InsertMany(context.TODO(), dataAry)
	if err != nil {
		return errors.Wrap(err, "problem inserting many patterns")
	}
	return nil
}

// *********************************************************
// Delete functions
// *********************************************************

func (repo *PatternRepo) DeleteAll() error {
	return repo.c.Drop(context.TODO())
}

// *********************************************************
// Find functions
// *********************************************************

func (repo *PatternRepo) FindOneAndReplace(pattern *model.Pattern) *model.Pattern {

	filter := bson.D{{"symbol", pattern.Symbol}, {"value", pattern.Value}}
	var update model.Pattern
	err := repo.c.FindOneAndReplace(context.TODO(), filter, pattern, repo.replaceOpt).Decode(&update)
	if err != nil {
		log.Warnf("problem replacing pattern due to %v", err)
	}
	return &update
}

func (repo *PatternRepo) FindAndReplace(pattern *model.Pattern) *model.Pattern {

	filter := bson.D{{"symbol", pattern.Symbol}, {"value", pattern.Value}}
	var update model.Pattern
	err := repo.c.FindOneAndReplace(context.TODO(), filter, pattern, repo.replaceOpt).Decode(&update)
	if err != nil {
		log.Warnf("problem replacing pattern due to %v", err)
	}
	return &update
}

func (repo *PatternRepo) FindBySymbol(symbol string) ([]*model.Pattern, error) {

	filter := bson.D{{"symbol", symbol}}
	var findData []*model.Pattern
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
		var doc model.Pattern
		err = cur.Decode(&doc)
		if err != nil {
			results = multierror.Append(results, err)
			continue
		}
		findData = append(findData, &doc)
	}
	return findData, results
}

func (repo *PatternRepo) FindOneBySymbolAndValue(symbol, value string) (*model.Pattern, error) {

	filter := bson.D{{"symbol", symbol}, {"value", value}}

	var pattern model.Pattern
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

func (repo *PatternRepo) FindHighestUpProbability(density model.PatternDensity) (*model.Pattern, error) {

	var pipeline mongo.Pipeline
	switch density {
	case model.PatternDensityLow:
		pipeline = patternAggregateUpMaxLowDensity
	case model.PatternDensityMedium:
		pipeline = patternAggregateUpMaxMediumDensity
	case model.PatternDensityHigh:
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

	var pattern model.Pattern
	// Only get the first one
	if cur.Next(context.TODO()) {
		err := cur.Decode(&pattern)
		if err != nil {
			return &pattern, err
		}
	}

	return &pattern, nil
}

func (repo *PatternRepo) FindHighestDownProbability(density model.PatternDensity) (*model.Pattern, error) {

	var pipeline mongo.Pipeline
	switch density {
	case model.PatternDensityLow:
		pipeline = patternAggregateDownMaxLowDensity
	case model.PatternDensityMedium:
		pipeline = patternAggregateDownMaxMediumDensity
	case model.PatternDensityHigh:
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

	var pattern model.Pattern
	// Only get the first one
	if cur.Next(context.TODO()) {
		err := cur.Decode(&pattern)
		if err != nil {
			return &pattern, err
		}
	}

	return &pattern, nil
}

func (repo *PatternRepo) FindHighestNoChangeProbability(density model.PatternDensity) (*model.Pattern, error) {

	var pipeline mongo.Pipeline
	switch density {
	case model.PatternDensityLow:
		pipeline = patternAggregateNoChangeMaxLowDensity
	case model.PatternDensityMedium:
		pipeline = patternAggregateNoChangeMaxMediumDensity
	case model.PatternDensityHigh:
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

	var pattern model.Pattern
	// Only get the first one
	if cur.Next(context.TODO()) {
		err := cur.Decode(&pattern)
		if err != nil {
			return &pattern, err
		}
	}

	return &pattern, nil
}

func (repo *PatternRepo) FindLowestUpProbability(density model.PatternDensity) (*model.Pattern, error) {

	var pipeline mongo.Pipeline
	switch density {
	case model.PatternDensityLow:
		pipeline = patternAggregateUpMinLowDensity
	case model.PatternDensityMedium:
		pipeline = patternAggregateUpMinMediumDensity
	case model.PatternDensityHigh:
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

	var pattern model.Pattern
	// Only get the first one
	if cur.Next(context.TODO()) {
		err := cur.Decode(&pattern)
		if err != nil {
			return &pattern, err
		}
	}

	return &pattern, nil
}

func (repo *PatternRepo) FindLowestDownProbability(density model.PatternDensity) (*model.Pattern, error) {

	var pipeline mongo.Pipeline
	switch density {
	case model.PatternDensityLow:
		pipeline = patternAggregateDownMinLowDensity
	case model.PatternDensityMedium:
		pipeline = patternAggregateDownMinMediumDensity
	case model.PatternDensityHigh:
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

	var pattern model.Pattern
	// Only get the first one
	if cur.Next(context.TODO()) {
		err := cur.Decode(&pattern)
		if err != nil {
			return &pattern, err
		}
	}

	return &pattern, nil
}

func (repo *PatternRepo) FindLowestNoChangeProbability(density model.PatternDensity) (*model.Pattern, error) {

	var pipeline mongo.Pipeline
	switch density {
	case model.PatternDensityLow:
		pipeline = patternAggregateNoChangeMinLowDensity
	case model.PatternDensityMedium:
		pipeline = patternAggregateNoChangeMinMediumDensity
	case model.PatternDensityHigh:
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

	var pattern model.Pattern
	// Only get the first one
	if cur.Next(context.TODO()) {
		err := cur.Decode(&pattern)
		if err != nil {
			return &pattern, err
		}
	}

	return &pattern, nil
}
