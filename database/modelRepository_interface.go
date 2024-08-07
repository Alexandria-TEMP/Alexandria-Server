package database

type ModelRepositoryInterface[T Model] interface {
	Create(object T) error
	GetByID(id uint) (T, error)
	Update(object T) (T, error)
	Delete(id uint) error
	// GetFields(wanted []interface{}) ([]interface{}, error)
	Query(conds ...interface{}) ([]T, error)
	QueryPaginated(page, size int, conds ...interface{}) ([]T, error)
}

//go:generate mockgen -package=mocks -source=./modelRepository_interface.go -destination=../mocks/modelRepository_mock.go
