package commands

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/kameshsampath/go-cqrs-demo/config"
	"github.com/kameshsampath/go-cqrs-demo/dao"
	"github.com/twmb/franz-go/pkg/kgo"
	"gorm.io/gorm"
)

func getKey(t *dao.Todo) string {
	return strconv.Itoa(int(t.ID))
}

func send(ctx context.Context, client *kgo.Client, eventType string, todo *dao.Todo) {
	tx := ctx.Value("TRANSACTION").(*gorm.DB)
	log := config.Log
	data := map[string]interface{}{
		"event_type": eventType,
		"todo":       todo,
	}
	log.Debugf("Data:%#v to topic", data)
	//Marshall to bytes
	tb, err := json.Marshal(data)
	if err != nil {
		log.Errorf("error marshalling data %s to topic, %s", todo, err)
		tx.Rollback()
	}
	//Send message to Redpanda Topic
	tctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	r := &kgo.Record{
		Key:   []byte(getKey(todo)),
		Value: tb,
	}
	client.Produce(tctx, r, func(r *kgo.Record, err error) {
		if err != nil {
			log.Errorf("error while sending %#v to topic,%s", data, err)
			//TODO rollback the transaction
			tx.Rollback()
		}
		tx.Commit()
		cancel()
		ctx.Done()
	})
}
