package interfaces

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

//go:generate mockgen -package=mocks -source=./memberService_interface.go -destination=../../mocks/memberService_mock.go

type MemberService interface {
	GetMember(memberID uint) (*models.Member, error)
	CreateMember(memberForm *forms.MemberCreationForm, userFields *models.ScientificFieldTagContainer) (string, string, *models.Member, error)
	DeleteMember(memberID uint) error
	GetAllMembers() ([]*models.MemberShortFormDTO, error)
	LogInMember(memberAuthForm *forms.MemberAuthForm) (*models.LoggedInMemberDTO, error)
	RefreshToken(form *forms.TokenRefreshForm) (*models.TokenPairDTO, error)
}
