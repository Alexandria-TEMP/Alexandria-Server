package services

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type ScientificFieldTagContainerService struct {
	ContainerRepository database.ModelRepositoryInterface[*models.ScientificFieldTagContainer]
}

func (containerService *ScientificFieldTagContainerService) GetScientificFieldTagContainer(containerID uint) (*models.ScientificFieldTagContainer, error) {
	return containerService.ContainerRepository.GetByID(containerID)
}
