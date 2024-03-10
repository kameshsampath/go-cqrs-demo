package commands

import (
	"github.com/kameshsampath/go-cqrs-demo/dao"
	"github.com/twmb/franz-go/pkg/kgo"
	"gorm.io/gorm"
)

// Command defines the methods for all the commands to be used within this application
type Command interface {
	SaveAndPublish() error
}

type baseModel struct {
	DB        *gorm.DB
	Client    *kgo.Client
	Todo      *dao.Todo
	EventType string
}
