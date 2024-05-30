package services

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/tags"
)

type MemberService struct {
	MemberRepository database.ModelRepository[*models.Member]
}

func (memberService *MemberService) GetMember(userID uint) (*models.Member, error) {
	//get member by this id
	member, err := memberService.MemberRepository.GetByID(userID)
	return member, err
}

func (memberService *MemberService) CreateMember(form *forms.MemberCreationForm, tags []*tags.ScientificFieldTag) *models.Member {
	// creating a member
	// := is declaration + assignment

	//for now no input sanitization for the strings - so first name, last name, email, institution, etc.
	//however have to get tags somehow

	//okay okay so for all form ids
	//should get tags and add them by ids
	  

	member := &models.Member{
		FirstName:   		form.FirstName,
		LastName:    		form.LastName,
		Email:       		form.Email,
		Password:    		form.Password,
		Institution: 		form.Institution,
		ScientificFieldTags: 	tags,
	}

	// TODO: add new member to repository

	return member
}

func (memberService *MemberService) UpdateMember(_ *models.Member) error {
	// TODO: database call to update
	return nil
}
