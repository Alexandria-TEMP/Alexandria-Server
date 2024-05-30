package services

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type MemberService struct {
	MemberRepository database.ModelRepository[*models.Member]
}

func (memberService *MemberService) GetMember(_ uint64) (*models.Member, error) {
	// TODO: database interaction?
	return new(models.Member), nil
}

func (memberService *MemberService) CreateMember(form *forms.MemberCreationForm) *models.Member {
	// creating a member
	// := is declaration + assignment
	member := &models.Member{
		FirstName:   form.FirstName,
		LastName:    form.LastName,
		Email:       form.Email,
		Password:    form.Password,
		Institution: form.Institution,
	}

	// TODO: add new member to repository

	return member
}

func (memberService *MemberService) UpdateMember(_ *models.Member) error {
	// TODO: database call to update
	return nil
}
