package dao

import (
	"fmt"

	"gorm.io/gorm"
)

// Todo represents the Todo entity int he Database
type Todo struct {
	gorm.Model
	Title       string `gorm:"column:title;not null" json:"title"`
	Category    string `gorm:"column:category;default:personal" json:"category,omitempty"`
	Description string `gorm:"column:description" json:"description,omitempty"`
	Status      bool   `gorm:"column:status;default:0" json:"status,omitempty"`
}

// TodoDAO defines the Data Access Methods for Todo entity
type TodoDAO interface {
	ListAll(db *gorm.DB) ([]Todo, error)
	ListByStatus(db *gorm.DB) ([]Todo, error)
	Save(db *gorm.DB) error
	Delete(db *gorm.DB) error
}

func (t Todo) String() string {
	return fmt.Sprintf(`Todo:{ ID:"%d", Title:"%s", Description:"%s", Category:"%s", Status:"%v" }`, t.ID, t.Title, t.Description, t.Category, t.Status)
}
