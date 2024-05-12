package interfaces

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/forms"
)

// run to create the mock
//go:generate mockgen -source=./userService_interface.go -destination=../../mocks/userService_mock.go

type UserService interface {
	GetMember(userID uint64) (*models.Member, error)
	CreateMember(memberForm *forms.MemberCreationForm) *models.Member
	UpdateMember(updatedMember *models.Member) error

	GetCollaborator(collaboratorID uint64) (*models.Collaborator, error)
	CreateCollaborator(collaboratorForm *forms.CollaboratorCreationForm) *models.Collaborator
	UpdateCollaborator(updatedCollaborator *models.Collaborator) error
}