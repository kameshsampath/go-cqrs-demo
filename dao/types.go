package dao

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

const MongoCollection = "todos"

// Todo represents the Todo entity int he Database
type Todo struct {
	gorm.Model
	Title       string `gorm:"column:title;not null" json:"title"`
	Category    string `gorm:"column:category;default:personal" json:"category,omitempty"`
	Description string `gorm:"column:description" json:"description,omitempty"`
	Status      bool   `gorm:"column:status;default:0" json:"status,omitempty"`
}

// BTodo represents the Todo entity BSON document
type BTodo struct {
	ObjectID    primitive.ObjectID `bson:"_id" json:"_id"`
	ID          uint               `bson:"id" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Category    string             `bson:"category" json:"category,omitempty"`
	Description string             `bson:"description" json:"description,omitempty"`
	Status      bool               `bson:"status" json:"status,omitempty"`
}

// Writer defines all the methods that will result in CUD of Todo
type Writer interface {
	Save(db *gorm.DB) error
	Delete(db *gorm.DB) error
}

// Reader defines all methods that will result in Read operations
type Reader interface {
	ListAll(col *mongo.Collection) ([]BTodo, error)
	ListByStatus(col *mongo.Collection) ([]BTodo, error)
	ListByCategory(col *mongo.Collection) ([]BTodo, error)
}

func (t Todo) String() string {
	return fmt.Sprintf(`Todo:{ ID:"%d", Title:"%s", Description:"%s", Category:"%s", Status:"%v" }`, t.ID, t.Title, t.Description, t.Category, t.Status)
}
