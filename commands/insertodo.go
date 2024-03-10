package commands

import (
	"github.com/kameshsampath/go-cqrs-demo/dao"
	"gorm.io/gorm"
)

// InsertTodo handles the event model for inserting todo
type InsertTodo struct {
	baseModel
}

func NewInsertCommand() *InsertTodo {
	return &InsertTodo{
		baseModel: baseModel{
			Todo:      &dao.Todo{},
			EventType: "insert",
		},
	}
}

// SaveAndPublish implements Command t
func (i *InsertTodo) SaveAndPublish() error {
	//run the entire stuff within a transaction
	//commit if all steps succeed or rollback
	err := i.DB.Transaction(func(tx *gorm.DB) error {
		// Insert TODO into the database
		if err := i.Todo.Save(tx); err != nil {
			return err
		}
		// send message to Redpanda Topic
		send(i.Client, i.EventType, i.Todo)
		return nil
	})
	return err
}

var _ Command = (*InsertTodo)(nil)
