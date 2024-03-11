package dao

import (
	"context"
	"time"

	"github.com/kameshsampath/go-cqrs-demo/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var log = config.Log

func W(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	//Migrate Schema
	err = db.AutoMigrate(&Todo{})

	return db, err
}

// Save implements Writer for handling create and update.
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

// Delete implements Writer.
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

var _ Writer = (*Todo)(nil)
