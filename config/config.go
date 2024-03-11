package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v10"
	"go.uber.org/zap"
)

type ConfigOpt func(*Config)

// Config sets the configuration for the Redpanda server
type Config struct {
	DBFile           string   `env:"DB_FILE" envDefault:"data/todo.db"`
	ConsumerGroupID  string   `env:"CONSUMER_GROUP_ID" envDefault:"todo-cqrs"`
	Seeds            []string `env:"RPK_BROKERS" envSeparator:"," envDefault:"localhost:19092"`
	Topics           []string `env:"TOPICS" envSeparator:"," envDefault:"todo"`
	AutoOffsetReset  string   `env:"AUTO_OFFSET_RESET" envSeparator:"," envDefault:"start"`
	DatabaseHost     string   `env:"PGHOST" envDefault:"localhost"`
	DatabasePort     string   `env:"PGPORT" envDefault:"5432"`
	DatabaseUser     string   `env:"PGUSER" envDefault:"demo"`
	DatabasePassword string   `env:"PGPASSWORD" envDefault:"superS3cret!"`
	DatabaseName     string   `env:"PGDATABASE" envDefault:"todos"`
	AtlasHost        string   `env:"ATLAS_HOST" envDefault:"localhost"`
	AtlasPort        string   `env:"ATLAS_PORT" envDefault:"27778"`
	AtlasUser        string   `env:"ATLAS_USER" envDefault:"demo"`
	AtlasPassword    string   `env:"ATLAS_PASSWORD" envDefault:"superS3cret!"`
}

var Log *zap.SugaredLogger

func init() {
	// Setup Log
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	Log = logger.Sugar()

}

func New(opts ...ConfigOpt) *Config {
	var config = new(Config)

	if err := env.Parse(config); err != nil {
		log.Fatalf("error parsing config, %v", err)
	}

	for _, opt := range opts {
		opt(config)
	}

	return config
}

func WithAtlasUser(atlasUser string) ConfigOpt {
	return func(c *Config) {
		c.AtlasUser = atlasUser
	}
}

func WithAtlasPassword(atlasPassword string) ConfigOpt {
	return func(c *Config) {
		c.AtlasPassword = atlasPassword
	}
}

func WithAtlasHost(atlasHost string) ConfigOpt {
	return func(c *Config) {
		c.AtlasHost = atlasHost
	}
}

func WithAtlasPort(atlasPort string) ConfigOpt {
	return func(c *Config) {
		c.AtlasPort = atlasPort
	}
}

func WithDatabaseUser(dbUser string) ConfigOpt {
	return func(c *Config) {
		c.DatabaseName = dbUser
	}
}

func WithDatabasePassword(dbUserPassword string) ConfigOpt {
	return func(c *Config) {
		c.DatabasePassword = dbUserPassword
	}
}

func WithDatabase(dbName string) ConfigOpt {
	return func(c *Config) {
		c.DatabaseName = dbName
	}
}
func WithDatabaseHost(dbHost string) ConfigOpt {
	return func(c *Config) {
		c.DatabaseHost = dbHost
	}
}
func WithDatabasePort(dbPort string) ConfigOpt {
	return func(c *Config) {
		c.DatabasePort = dbPort
	}
}

func WithTopics(topics []string) ConfigOpt {
	return func(c *Config) {
		c.Topics = topics
	}
}

func WithConsumerGroupID(groupID string) ConfigOpt {
	return func(c *Config) {
		c.ConsumerGroupID = groupID
	}
}

func WithSeeds(seeds []string) ConfigOpt {
	return func(c *Config) {
		c.Seeds = seeds
	}
}

func WithAutoOffsetReset(offsetReset string) ConfigOpt {
	return func(c *Config) {
		c.AutoOffsetReset = offsetReset
	}
}

// DefaultProducerTopic gets the default topic that will be used as the producer topic
func (c *Config) DefaultProducerTopic() string {
	return c.Topics[0]
}

// DSN returns the Database DSN
func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DatabaseHost,
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseName,
		c.DatabasePort)
}

// MongoURI returns the MongoDB Atlas connection URI
func (c *Config) MongoURI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s/?directConnection=true",
		c.AtlasUser,
		c.AtlasPassword,
		c.AtlasHost,
		c.AtlasPort)
}
