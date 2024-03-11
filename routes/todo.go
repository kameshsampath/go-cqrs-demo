package routes

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/kameshsampath/go-cqrs-demo/commands"
	"github.com/kameshsampath/go-cqrs-demo/dao"
	"github.com/labstack/echo/v4"
	"github.com/twmb/franz-go/pkg/kgo"
	"gorm.io/gorm"
)

var lock sync.Mutex

func ListAll(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	t := &dao.Todo{}
	db := c.Get("DB").(*gorm.DB)
	todos, err := t.ListAll(db)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, todos)
}

func ListByStatus(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	t := &dao.Todo{Status: to_bool(c.Param("status"))}
	db := c.Get("DB").(*gorm.DB)
	todos, err := t.ListByStatus(db)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, todos)
}

func Create(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	// Create InsertTodo command
	ic := commands.NewInsertCommand()
	ic.DB = c.Get("DB").(*gorm.DB)
	ic.Client = c.Get("RP_CLIENT").(*kgo.Client)
	// Bind the request body
	if err := c.Bind(ic.Todo); err != nil {
		return err
	}
	// Save the Todo to DB and Publish the data to Redpanda topic
	err := ic.SaveAndPublish()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, ic.Todo)
}

func Update(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	// Create UpdateTodo command
	uc := commands.NewUpdateCommand()
	uc.DB = c.Get("DB").(*gorm.DB)
	uc.Client = c.Get("RP_CLIENT").(*kgo.Client)
	err := (&echo.DefaultBinder{}).BindBody(c, uc.Todo)
	if err != nil {
		return err
	}
	// Save the Todo to DB and Publish the data to Redpanda topic
	err = uc.SaveAndPublish()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusNoContent, nil)
}

func Delete(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	// Create DeleteTodo command
	dc := commands.NewDeleteCommand()
	dc.DB = c.Get("DB").(*gorm.DB)
	dc.Client = c.Get("RP_CLIENT").(*kgo.Client)
	id, _ := strconv.Atoi(c.Param("id"))
	dc.Todo.ID = uint(id)
	// Save the Todo to DB and Publish the data to Redpanda topic
	err := dc.SaveAndPublish()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusNoContent, nil)
}

func to_bool(str string) bool {
	switch str {
	case "yes", "done", "ok", "true":
		return true
	default:
		return false
	}
}
