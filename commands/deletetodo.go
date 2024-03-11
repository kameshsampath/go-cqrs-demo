package commands

import (
	"context"

	"github.com/kameshsampath/go-cqrs-demo/dao"
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
	tx := d.DB.Begin()
	// Delete TODO from the database
	if err := d.Todo.Delete(tx); err != nil {
		tx.Rollback()
		return err
	}
	// if the message sending to topic encounters error
	// rollback the transaction
	ctx := context.WithValue(context.TODO(), "TRANSACTION", tx)
	go send(ctx, d.Client, d.EventType, d.Todo)

	return nil
}

var _ Command = (*DeleteTodo)(nil)
