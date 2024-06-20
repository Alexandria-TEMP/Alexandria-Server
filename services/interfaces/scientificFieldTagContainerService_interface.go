package interfaces

import "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"

//go:generate mockgen -package=mocks -source=./scientificFieldTagContainerService_interface.go -destination=../../mocks/scientificFieldTagContainerService_mock.go

type ScientificFieldTagContainerService interface {
	// GetScientificFieldTagContainer gets a ScientificFieldTagContainer by ID from the database
	GetScientificFieldTagContainer(containerID uint) (*models.ScientificFieldTagContainer, error)
}
