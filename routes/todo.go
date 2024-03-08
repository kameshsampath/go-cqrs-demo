package routes

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/kameshsampath/go-cqrs-demo/dao"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var lock sync.Mutex
var log *zap.SugaredLogger

func init() {
	// Setup Logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	log = logger.Sugar()
}

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
	t := &dao.Todo{}
	if err := c.Bind(t); err != nil {
		return err
	}
	db := c.Get("DB").(*gorm.DB)
	err := t.Save(db)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, t)
}

func Update(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	t := dao.Todo{}
	err := (&echo.DefaultBinder{}).BindBody(c, &t)
	if err != nil {
		return err
	}
	log.Debugf("Update:%s", t)
	db := c.Get("DB").(*gorm.DB)
	err = t.Save(db)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusNoContent, "")
}

func Delete(c echo.Context) error {
	lock.Lock()
	defer lock.Unlock()
	t := dao.Todo{}
	id, _ := strconv.Atoi(c.Param("id"))
	t.ID = uint(id)
	db := c.Get("DB").(*gorm.DB)
	err := t.Delete(db)
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
