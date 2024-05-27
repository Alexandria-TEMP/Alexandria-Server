package interfaces

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

// run to create the mock
//go:generate mockgen -package=mocks -source=./memberService_interface.go -destination=../../mocks/memberService_mock.go

type MemberService interface {
	GetMember(userID uint64) (*models.Member, error)
	CreateMember(memberForm *forms.MemberCreationForm) *models.Member
	UpdateMember(updatedMember *models.Member) error

	GetCollaborator(collaboratorID uint64) (*models.PostCollaborator, error)
	CreateCollaborator(collaboratorForm *forms.CollaboratorCreationForm) *models.PostCollaborator
	UpdateCollaborator(updatedCollaborator *models.PostCollaborator) error
}
