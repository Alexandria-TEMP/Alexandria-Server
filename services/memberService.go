package services

import (
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type MemberService struct {
	// dont know how to do database connections?
	// but i think they go in here
}

func (memberService *MemberService) GetMember(_ uint) (*models.Member, error) {
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

func (memberService *MemberService) GetCollaborator(_ uint) (*models.PostCollaborator, error) {
	// TODO: actually get from database based on UUID
	return new(models.PostCollaborator), nil
}

func (memberService *MemberService) CreateCollaborator(form *forms.CollaboratorCreationForm) *models.PostCollaborator {
	collaborator := &models.PostCollaborator{
		Member: form.Member,
		// TODO: is this correct? will it assign the right thing?
		// honestly have no clue yet
		CollaborationType: models.CollaborationType(form.CollaborationType),
	}

	// TODO: add this one to the database as well

	return collaborator
}

func (memberService *MemberService) UpdateCollaborator(_ *models.PostCollaborator) error {
	// TODO: update data in database
	return nil
}
