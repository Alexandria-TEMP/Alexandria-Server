package services

import (
	rand "math/rand"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models/forms"
)

type UserService struct {
	// dont know how to do database connections?
	//but i think they go in here
}

func (userService *UserService) GetMember(_ uint64) (*models.Member, error) {
	//TODO: database interaction?
	return new(models.Member), nil
}

func (userService *UserService) CreateMember(form *forms.MemberCreationForm) *models.Member {
	//creating a member
	//:= is declaration + assignment
	member := &models.Member{
		UserID:      rand.Uint64(),
		FirstName:   form.FirstName,
		LastName:    form.LastName,
		Email:       form.Email,
		Password:    form.Password,
		Institution: form.Institution,
		Posts:       form.Posts,
		Discussions: form.Discussions,
		Reviews:     form.Reviews,
	}

	//TODO: add new member to repository

	return member

}
func (userService *UserService) UpdateMember (_ *models.Member) error {
	//TODO: database call to update 
	return nil
}
func (userService *UserService) GetCollaborator(_ uint64) (*models.Collaborator, error) {
	//TODO: actually get from database based on UUID
	return new(models.Collaborator), nil
}

func (userService *UserService) CreateCollaborator(form *forms.CollaboratorCreationForm) *models.Collaborator {
	collaborator := &models.Collaborator{
		CollaboratorID:		rand.Uint64(),
		Member:				form.Member,
		//TODO: is this correct? will it assign the right thing?
		//honestly have no clue yet
		CollaborationType:	models.CollaborationType(form.CollaborationType),
	}

	//TODO: add this one to the database as well

	return collaborator
}

func (userService *UserService) UpdateCollaborator(_ *models.Collaborator) error {
	//TODO: update data in database

	return nil
}
