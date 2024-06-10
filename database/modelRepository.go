package database

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Performs CRUD operations on a given type of database model to the database.
// Type T must be a pointer to a struct, e.g. *Member.
// Example usage: repo := ModelRepository[*Member] { ... }
type ModelRepository[T Model] struct {
	Database *gorm.DB
}

// Create an object in the database. The passed object is updated with a new ID.
//
// T must be a pointer to a Model type (as outlined above).
// In general, you should *not* specify the ID (leave it blank during struct creation)
//
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
	result := repo.Database.Preload(clause.Associations).First(&found, id)

	if result.Error != nil {
		var zero T
		return zero, result.Error
	}

	return found, nil
}

func (repo *ModelRepository[T]) GetBy(match T) ([]T, error) {
	var found []T
	result := repo.Database.Where(match).Find(&found)

	if result.Error != nil {
		var zero []T
		return zero, result.Error
	}

	return found, nil
}

func (repo *ModelRepository[T]) Update(object T) (T, error) {
	// Ensure a model with this ID already exists
	id := object.GetID()

	result := repo.Database.First(new(T), id)
	if result.Error != nil {
		var zero T
		return zero, fmt.Errorf("could not find model with ID %d to update: %w", id, result.Error)
	}

	// Save the new data
	result = repo.Database.Save(object)
	if result.Error != nil {
		var zero T
		return zero, fmt.Errorf("could not update model with ID %d: %w", id, result.Error)
	}

	return repo.GetByID(id)
}

func (repo *ModelRepository[T]) Delete(id uint) error {
	// Ensure a model with this ID already exists
	result := repo.Database.First(new(T), id)
	if result.Error != nil {
		return fmt.Errorf("could not find model with ID %d to delete: %w", id, result.Error)
	}

	result = repo.Database.Delete(new(T), id)

	return result.Error
}
