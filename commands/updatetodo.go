package commands

import (
	"github.com/kameshsampath/go-cqrs-demo/dao"
	"gorm.io/gorm"
)

// UpdateTodo handles the event model for updating todo
type UpdateTodo struct {
	baseModel
}

// NewUpdateCommand creates a new instance of the Update Todo command
func NewUpdateCommand() *UpdateTodo {
	return &UpdateTodo{
		baseModel: baseModel{
			Todo:      &dao.Todo{},
			EventType: "update",
		},
	}
}

// SaveAndPublish implements Command
func (u *UpdateTodo) SaveAndPublish() error {
	//run the entire stuff within a transaction
	//commit if all steps succeed or rollback
	err := u.DB.Transaction(func(tx *gorm.DB) error {
		// Update TODO into the database
		if err := u.Todo.Save(tx); err != nil {
			return err
		}
		// send message to Redpanda Topic
		send(u.Client, u.EventType, u.Todo)
		return nil
	})
	return err
}

var _ Command = (*UpdateTodo)(nil)
