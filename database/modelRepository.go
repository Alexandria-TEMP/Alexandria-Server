package database

import (
	"log"

	"gorm.io/gorm"
)

// Interface that any database model must adhere to. Used by the ModelRepository.
type Model interface {
	GetID() uint
}

// Performs CRUD operations on a given type of database model to the database.
// Type T must be a pointer to a struct, e.g. *Member.
// Example usage: repo := ModelRepository[*Member] { db: ... }
type ModelRepository[T Model] struct {
	Database *gorm.DB
}

// Create an object in the database. The passed object's initial ID field is ignored,
// and a fresh new ID will be assigned to it. This modifies the original object, as
// T must be a pointer to a Model type (as outlined above).
func (repo *ModelRepository[T]) Create(object T) error {
	result := repo.Database.Create(object)

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
