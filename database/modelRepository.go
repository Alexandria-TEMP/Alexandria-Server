package database

import (
	"log"

	"gorm.io/gorm"
)

// Interface that any database model must adhere to. Used by the ModelRepository.
type Model interface {
	getID() uint
}

// Performs CRUD operations on a given type of database model to the database.
type ModelRepository[T Model] struct {
	db *gorm.DB
}

// Create an object in the database. The passed object's initial ID field is ignored,
// and a fresh new ID will be assigned to it. This modifies the original object.
func (repo *ModelRepository[T]) Create(object *T) error {
	result := repo.db.Create(object)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ModelRepository[T]) GetByID(id uint) {
	log.Fatal("TODO GetByID")
}

func (repo *ModelRepository[T]) Update(object T) {
	log.Fatal("TODO Update")
}

func (repo *ModelRepository[T]) Delete(object T) {
	log.Fatal("TODO Delete")
}

func (repo *ModelRepository[T]) DeleteByID(ID uint) {
	log.Fatal("TODO Delete")
}
