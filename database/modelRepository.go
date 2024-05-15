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
// Example usage: repo := ModelRepository[*Member] { ... }
type ModelRepository[T Model] struct {
	Database *gorm.DB
}

// Create an object in the database.
// T must be a pointer to a Model type (as outlined above).
// In general, you should *not* specify the ID (leave it blank during struct creation
// If the ID is specified:
// - if an object with the given ID already exists, errors.
// - otherwise, creates the object with that ID.
func (repo *ModelRepository[T]) Create(object T) error {
	result := repo.Database.Create(object)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ModelRepository[T]) GetByID(id uint) (T, error) {
	var found T
	result := repo.Database.First(&found, id)

	if result.Error != nil {
		var zero T
		return zero, result.Error
	}

	return found, nil
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
