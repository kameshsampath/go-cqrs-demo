package commands

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/kameshsampath/go-cqrs-demo/config"
	"github.com/kameshsampath/go-cqrs-demo/dao"
	"github.com/twmb/franz-go/pkg/kgo"
)

func getKey(t *dao.Todo) string {
	return strconv.Itoa(int(t.ID))
}

func send(client *kgo.Client, eventType string, todo *dao.Todo) {
	ctx := context.TODO()
	go func(ctx context.Context) {
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
		}
		//Send message to Redpanda Topic
		tctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()
		r := &kgo.Record{
			Key:   []byte(getKey(todo)),
			Value: tb,
		}
		if err := client.ProduceSync(tctx, r).FirstErr(); err != nil {
			log.Errorf("error sending %#v to topic, %s", data, err)
		}
		ctx.Done()
	}(ctx)
}
