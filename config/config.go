package config

import (
	"log"

	"github.com/caarlos0/env/v10"
	"go.uber.org/zap"
)

// Config sets the configuration for the Redpanda server
type Config struct {
	DBFile          string   `env:"DB_FILE" envDefault:"data/todo.db"`
	ConsumerGroupID string   `env:"CONSUMER_GROUP_ID" envDefault:"todo-cqrs"`
	Seeds           []string `env:"RPK_BROKERS" envSeparator:"," envDefault:"localhost:19092"`
	Topics          []string `env:"TOPICS" envSeparator:"," envDefault:"todo"`
}

var Log *zap.SugaredLogger

func init() {
	// Setup Log
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	Log = logger.Sugar()

}

func New() *Config {
	var config = new(Config)

	if err := env.Parse(config); err != nil {
		log.Fatalf("error parsing config, %v", err)
	}

	return config
}

// DefaultProducerTopic gets the default topic that will be used as the producer topic
func (c *Config) DefaultProducerTopic() string {
	return c.Topics[0]
}
