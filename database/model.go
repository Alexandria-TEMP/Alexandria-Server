package database

// Interface that any database model must adhere to. Used by the ModelRepository.
type Model interface {
	GetID() uint
}
