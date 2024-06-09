package database

type ModelRepositoryInterface[T Model] interface {
	Create(object T) error
	GetByID(id uint) (T, error)
	GetBy(match T) ([]T, error)
	Update(object T) (T, error)
	Delete(id uint) error
}

//go:generate mockgen -package=mocks -source=./modelRepository_interface.go -destination=../mocks/modelRepository_mock.go
