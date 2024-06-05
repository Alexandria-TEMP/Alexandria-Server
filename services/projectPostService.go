package services

import (
	"fmt"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type ProjectPostService struct {
}

func (projectPostService *ProjectPostService) GetProjectPost(_ uint) (*models.ProjectPost, error) {
	return nil, fmt.Errorf("TODO")
}

func (projectPostService *ProjectPostService) CreateProjectPost(_ *forms.ProjectPostCreationForm) (*models.ProjectPost, error) {
	return nil, fmt.Errorf("TODO")
}

func (projectPostService *ProjectPostService) UpdateProjectPost(_ *models.ProjectPost) error {
	return fmt.Errorf("TODO")
}
