package services

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type MemberService struct {
	MemberRepository database.ModelRepositoryInterface[*models.Member]
}

func (memberService *MemberService) GetMember(memberID uint) (*models.Member, error) {
	// get member by this id
	member, err := memberService.MemberRepository.GetByID(memberID)
	return member, err
}

func (memberService *MemberService) CreateMember(form *forms.MemberCreationForm, userFields *models.ScientificFieldTagContainer) (*models.Member, error) {
	// for now no input sanitization for the strings - so first name, last name, email, institution, etc.
	// however have to get tags somehow
	member := &models.Member{
		FirstName:                   form.FirstName,
		LastName:                    form.LastName,
		Email:                       form.Email,
		Password:                    form.Password,
		Institution:                 form.Institution,
		ScientificFieldTagContainer: *userFields,
	}

	err := memberService.MemberRepository.Create(member)
	if err != nil {
		return nil, err
	}

	return member, err
}

func (memberService *MemberService) DeleteMember(memberID uint) error {
	err := memberService.MemberRepository.Delete(memberID)
	return err
}

func (memberService *MemberService) GetAllMembers() ([]*models.MemberShortFormDTO, error) {
	members, err := memberService.MemberRepository.Query()

	shortFormDTOs := make([]*models.MemberShortFormDTO, len(members))
	for i, member := range members {
		shortFormDTOs[i] = &models.MemberShortFormDTO{
			ID:        member.ID,
			FirstName: member.FirstName,
			LastName:  member.LastName,
		}
	}

	return shortFormDTOs, err
}
