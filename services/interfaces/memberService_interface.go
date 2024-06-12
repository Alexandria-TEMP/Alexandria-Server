package interfaces

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
)

//go:generate mockgen -package=mocks -source=./memberService_interface.go -destination=../../mocks/memberService_mock.go

type MemberService interface {
	GetMember(memberID uint) (*models.Member, error)
	CreateMember(memberForm *forms.MemberCreationForm, userFields *tags.ScientificFieldTagContainer) (*models.Member, error)
	UpdateMember(updatedMember *models.MemberDTO, userFields *tags.ScientificFieldTagContainer) error
	DeleteMember(memberID uint) error
	GetAllMembers() ([]*models.MemberShortFormDTO, error)
}
