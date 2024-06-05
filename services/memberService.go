package services

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
)

type MemberService struct {
	MemberRepository database.RepositoryInterface[*models.Member]
}

func (memberService *MemberService) GetMember(userID uint) (*models.Member, error) {
	// get member by this id
	member, err := memberService.MemberRepository.GetByID(userID)
	return member, err
}

func (memberService *MemberService) CreateMember(form *forms.MemberCreationForm, userFields tags.ScientificFieldTagContainer) (*models.Member, error) {
	// for now no input sanitization for the strings - so first name, last name, email, institution, etc.
	// however have to get tags somehow
	member := &models.Member{
		FirstName:                   form.FirstName,
		LastName:                    form.LastName,
		Email:                       form.Email,
		Password:                    form.Password,
		Institution:                 form.Institution,
		ScientificFieldTagContainer: userFields,
	}

	err := memberService.MemberRepository.Create(member)
	if err != nil {
		return nil, err
	}

	return member, err
}

func (memberService *MemberService) UpdateMember(memberDTO *models.MemberDTO, userFields tags.ScientificFieldTagContainer) error {
	member := &models.Member{
		FirstName:                   memberDTO.FirstName,
		LastName:                    memberDTO.LastName,
		Email:                       memberDTO.Email,
		Password:                    memberDTO.Password,
		Institution:                 memberDTO.Institution,
		ScientificFieldTagContainer: userFields,
	}
	_, err := memberService.MemberRepository.Update(member)
	return err
}

func (memberService *MemberService) DeleteMember(userID uint) error {
	err := memberService.MemberRepository.Delete(userID)
	return err
}

func (memberService *MemberService) GetAllMembers() ([]uint, error) {
	result, err := memberService.MemberRepository.GetAllIDs()
	return result, err
}
