package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/kameshsampath/go-cqrs-demo/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var log = config.Log

func New(dbFile string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("file:%s", dbFile)
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	//Migrate Schema
	err = db.AutoMigrate(&Todo{})

	return db, err
}

// Save implements TodoDAO for handling create and update.
func (t *Todo) Save(db *gorm.DB) error {
	log.Debugf("Saving %s", t)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	tx := db.WithContext(ctx).Save(t)
	if err := tx.Error; err != nil {
		log.Errorf("Error adding todo %s, %s", t, err)
		return err
	}
	return nil
}

// Delete implements TodoDAO.
func (t *Todo) Delete(db *gorm.DB) error {
	log.Debugf("Deleting %s", t)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	tx := db.WithContext(ctx).Delete(t)
	if err := tx.Error; err != nil {
		log.Errorf("Error deleting todo %s, %s", t, err)
		return err
	}
	return nil
}

// ListAll implements TodoDAO.
func (t *Todo) ListAll(db *gorm.DB) (todos []Todo, err error) {
	log.Debugf("Getting all todos")
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	tx := db.WithContext(ctx).Find(&todos)
	if err = tx.Error; err != nil {
		return nil, err
	}
	log.Debugf("Retrieved %d todos", tx.RowsAffected)
	return todos, err
}

// ListByStatus implements TodoDAO.
func (t *Todo) ListByStatus(db *gorm.DB) (todos []Todo, err error) {
	log.Debugf("Getting todos with status %v", t.Status)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	tx := db.WithContext(ctx).Where("status = ?", t.Status).Find(&todos)
	if err = tx.Error; err != nil {
		return nil, err
	}
	log.Debugf("Retrieved %d todos for status %v", tx.RowsAffected, t.Status)
	return todos, err
}

var _ TodoDAO = (*Todo)(nil)
