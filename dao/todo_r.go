package dao

import (
	"context"

	"github.com/kameshsampath/go-cqrs-demo/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func R(cfg *config.Config) (*mongo.Collection, error) {
	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.MongoURI()))
	if err != nil {
		log.Fatalf("error connecting to mongodb, %s", err)
	}
	//Check connectivity
	if err := db.Ping(context.TODO(), nil); err != nil {
		log.Fatalf("error connecting to mongodb, %s", err)
	}

	col := db.
		Database(cfg.AtlasDatabase).
		Collection(MongoCollection)

	return col, err
}

// ListAll implements Reader.
func (t *Todo) ListAll(col *mongo.Collection) ([]BTodo, error) {
	ctx := context.TODO()
	cur, err := col.Find(ctx, bson.D{})
	if err != nil {
		log.Errorf("error getting all todos,%s", err)
		return nil, err
	}
	return buildResults(ctx, cur)
}

// ListByCategory implements Reader.
func (t *Todo) ListByCategory(col *mongo.Collection) ([]BTodo, error) {
	ctx := context.TODO()
	filter := bson.D{{
		Key:   "category",
		Value: bson.D{{Key: "$eq", Value: t.Category}},
	}}
	cur, err := col.Find(ctx, filter)
	if err != nil {
		log.Errorf("error getting todos by category %s,%s", t.Category, err)
		return nil, err
	}
	return buildResults(ctx, cur)
}

// ListByStatus implements Reader.
func (t *Todo) ListByStatus(col *mongo.Collection) ([]BTodo, error) {
	ctx := context.TODO()
	filter := bson.D{{
		Key:   "status",
		Value: bson.D{{Key: "$eq", Value: t.Status}},
	}}
	cur, err := col.Find(ctx, filter)
	if err != nil {
		log.Errorf("error getting todos by status %s,%s", t.Status, err)
		return nil, err
	}
	return buildResults(ctx, cur)
}

func buildResults(ctx context.Context, cur *mongo.Cursor) ([]BTodo, error) {
	var todos []BTodo
	err := cur.All(ctx, &todos)
	if err != nil {
		log.Errorf("error getting todos,%s", err)
		return nil, err
	}
	return todos, nil
}

var _ Reader = (*Todo)(nil)
