package utils

import (
	"github.com/kameshsampath/go-cqrs-demo/config"
	"github.com/twmb/franz-go/pkg/kgo"
)

func NewClient(cfg *config.Config) (*kgo.Client, error) {
	if cfg == nil {
		cfg = config.New()
	}
	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Seeds...),
		kgo.ConsumeTopics(cfg.Topics...),
		kgo.DefaultProduceTopic(cfg.DefaultProducerTopic()),
		kgo.ConsumerGroup(cfg.ConsumerGroupID),
		kgo.AllowAutoTopicCreation(),
		kgo.DisableAutoCommit(),
		//TODO use config value
		kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()),
	)
	if err != nil {
		return nil, err
	}

	return client, err
}
