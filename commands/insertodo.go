package commands

import (
	"context"

	"github.com/kameshsampath/go-cqrs-demo/dao"
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
func (i *InsertTodo) SaveAndPublish() (err error) {
	//run the entire stuff within a transaction
	//commit if all steps succeed or rollback
	tx := i.DB.Begin()
	// Insert TODO into the database
	if err = i.Todo.Save(tx); err != nil {
		tx.Rollback()
		return err
	}
	// if the message sending to topic encounters error
	// rollback the transaction
	ctx := context.WithValue(context.TODO(), "TRANSACTION", tx)
	go send(ctx, i.Client, i.EventType, i.Todo)

	return err
}

var _ Command = (*InsertTodo)(nil)
