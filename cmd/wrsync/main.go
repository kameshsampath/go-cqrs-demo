package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/kameshsampath/go-cqrs-demo/config"
	"github.com/kameshsampath/go-cqrs-demo/dao"
	"github.com/kameshsampath/go-cqrs-demo/utils"
	"github.com/labstack/gommon/log"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type eventData struct {
	EventType string   `json:"event_type"`
	Todo      dao.Todo `json:"todo"`
}

const mongoCollection = "todos"

var (
	db  *mongo.Client
	cfg *config.Config
)

func main() {
	log := config.Log
	//Setup the Redpand Client
	//TODO remove after tests
	groupID := fmt.Sprintf("group-id-%d", time.Now().UnixMilli())
	cfg = config.New(config.WithConsumerGroupID(groupID))
	log.Debugf("Config:%#v", cfg)
	client, err := utils.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	//Connect to MongoDB
	db, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.MongoURI()))
	if err != nil {
		log.Fatalf("error connecting to MongoDB, %s", err)
	}
	if err := db.Ping(context.TODO(), nil); err != nil {
		log.Fatalf("error connecting to MongoDB, %s", err)
	}

	// handle interrupts
	ctx, quit := signal.NotifyContext(context.Background(), os.Interrupt)
	defer quit()
poll:
	for {
		select {
		case <-ctx.Done():
			break poll
		default:
			//Process Records, since it will have a DB insert
			//putting it into its own goroutine for better concurrency
			go func(ctx context.Context) {
				//Poll records
				fetches := client.PollFetches(ctx)
				//log errors
				if errs := fetches.Errors(); len(errs) > 0 {
					log.Debugf("errors during fetch,%v", errs)
					for err := range errs {
						log.Error(err)
					}
				}
				fetches.EachPartition(func(p kgo.FetchTopicPartition) {
					log.Debugf("Processing partition,%d", p.Partition)
					p.EachRecord(processRecord)
				})
			}(ctx)
		}
	}
}

func processRecord(r *kgo.Record) {
	ctx := context.TODO()
	log.Debugf("Processing Record,%#v", r)
	log := config.Log
	v := r.Value
	data := new(eventData)
	if err := json.Unmarshal(v, &data); err == nil {
		log.Debugf("Data:%#v", data)
	}
	col := db.
		Database(cfg.AtlasDatabase).
		Collection(mongoCollection)
	switch data.EventType {
	case "insert", "update":
		//data that will be updated or inserted
		update := bson.D{
			{
				Key: "$set",
				Value: bson.D{
					{Key: "id", Value: data.Todo.ID},
					{Key: "title", Value: data.Todo.Title},
					{Key: "description", Value: data.Todo.Description},
					{Key: "category", Value: data.Todo.Category},
					{Key: "status", Value: data.Todo.Status},
					{Key: "createdAt", Value: data.Todo.CreatedAt},
					{Key: "deletedAt", Value: data.Todo.DeletedAt},
					{Key: "updatedAt", Value: data.Todo.UpdatedAt},
				},
			},
		}
		log.Infof("Saving document,%#v", update)
		//filter to use to find the matching records
		filter := bson.D{
			{Key: "id", Value: data.Todo.ID},
			{Key: "title", Value: data.Todo.Title},
		}
		//insert if it does not exist otherwise update the existing record
		opts := options.Update().SetUpsert(true)
		r, err := col.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Errorf("error saving %s, %s", update, err)
		}
		if r.UpsertedID == nil {
			log.Debugf("Updated %d document(s).", r.MatchedCount)
		} else {
			log.Debugf("Created document with ID:%s.", r.UpsertedID)
		}

	case "delete":
		// filter to find matching records to delete
		filter := bson.D{
			{Key: "id", Value: data.Todo.ID},
			{Key: "title", Value: data.Todo.Title},
		}
		log.Infof("Delete the record,%s", filter)
		r, err := col.DeleteOne(ctx, filter)
		if err != nil {
			log.Errorf("error deleting %s, %s", filter, err)
		}
		log.Debugf("Deleted %d document", r.DeletedCount)
	}
}
