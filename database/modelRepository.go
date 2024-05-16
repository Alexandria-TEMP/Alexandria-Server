package database

import (
	"fmt"

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
	database *gorm.DB
}

// Create an object in the database.
// T must be a pointer to a Model type (as outlined above).
// In general, you should *not* specify the ID (leave it blank during struct creation)
// If the ID is specified:
// - if an object with the given ID already exists, errors.
// - otherwise, creates the object with that ID.
func (repo *ModelRepository[T]) Create(object T) error {
	result := repo.database.Create(object)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ModelRepository[T]) GetByID(id uint) (T, error) {
	var found T
	result := repo.database.First(&found, id)

	if result.Error != nil {
		var zero T
		return zero, result.Error
	}

	return found, nil
}

func (repo *ModelRepository[T]) Update(object T) (T, error) {
	// Ensure a model with this ID already exists
	id := object.GetID()

	result := repo.database.First(new(T), id)
	if result.Error != nil {
		var zero T
		return zero, fmt.Errorf("could not find model with ID %d to update: %w", id, result.Error)
	}

	// Save the new data
	result = repo.database.Save(object)
	if result.Error != nil {
		var zero T
		return zero, fmt.Errorf("could not update model with ID %d: %w", id, result.Error)
	}

	// Return the newly saved object, because some of its fields
	// (e.g. last-updated) may have been changed automatically.
	return repo.GetByID(id)
}

func (repo *ModelRepository[T]) Delete(id uint) error {
	// Ensure a model with this ID already exists
	result := repo.database.First(new(T), id)
	if result.Error != nil {
		return fmt.Errorf("could not find model with ID %d to delete: %w", id, result.Error)
	}

	result = repo.database.Delete(new(T), id)

	return result.Error
}
