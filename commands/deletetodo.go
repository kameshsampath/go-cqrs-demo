package commands

import (
	"github.com/kameshsampath/go-cqrs-demo/dao"
	"gorm.io/gorm"
)

// DeleteTodo handles the event model for deleting a todo
type DeleteTodo struct {
	baseModel
}

// NewDeleteCommand creates a new Delete Todo Command
func NewDeleteCommand() *DeleteTodo {
	return &DeleteTodo{
		baseModel: baseModel{
			Todo:      &dao.Todo{},
			EventType: "delete",
		},
	}
}

// SaveAndPublish implements
func (d *DeleteTodo) SaveAndPublish() error {
	//run the entire stuff within a transaction
	//commit if all steps succeed or rollback
	err := d.DB.Transaction(func(tx *gorm.DB) error {
		// Delete TODO from the database
		if err := d.Todo.Delete(tx); err != nil {
			return err
		}
		// send message to Redpanda Topic
		send(d.Client, d.EventType, d.Todo)
		return nil
	})
	return err
}

var _ Command = (*DeleteTodo)(nil)
