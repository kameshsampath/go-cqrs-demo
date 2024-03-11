package commands

import (
	"context"

	"github.com/kameshsampath/go-cqrs-demo/dao"
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
	tx := u.DB.Begin()
	// Update TODO into the database
	if err := u.Todo.Save(tx); err != nil {
		tx.Rollback()
		return err
	}
	// if the message sending to topic encounters error
	// rollback the transaction
	ctx := context.WithValue(context.TODO(), "TRANSACTION", tx)
	go send(ctx, u.Client, u.EventType, u.Todo)

	return nil
}

var _ Command = (*UpdateTodo)(nil)
