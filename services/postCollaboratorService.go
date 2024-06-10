package services

import (
	"fmt"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type PostCollaboratorService struct {
	PostCollaboratorRepository database.ModelRepositoryInterface[*models.PostCollaborator]
	MemberRepository           database.ModelRepositoryInterface[*models.Member]
}

func (postCollaboratorService *PostCollaboratorService) GetPostCollaborator(id uint) (*models.PostCollaborator, error) {
	return postCollaboratorService.PostCollaboratorRepository.GetByID(id)
}

func (postCollaboratorService *PostCollaboratorService) MembersToPostCollaborators(memberIDs []uint, anonymous bool, collaborationType models.CollaborationType) ([]*models.PostCollaborator, error) {
	// If the list is anonymous, immediately return empty
	if anonymous {
		return []*models.PostCollaborator{}, nil
	}

	// If the list is not anonymous, check it has at least one author
	if len(memberIDs) < 1 {
		return []*models.PostCollaborator{}, fmt.Errorf("could not create post collaborators: must have at least one member")
	}

	postCollaborators := make([]*models.PostCollaborator, len(memberIDs))

	for i, memberID := range memberIDs {
		// Fetch the member from the database
		member, err := postCollaboratorService.MemberRepository.GetByID(memberID)
		if err != nil {
			return nil, fmt.Errorf("could not create post collaborators: %w", err)
		}

		newPostCollaborator := models.PostCollaborator{
			Member:            *member,
			CollaborationType: collaborationType,
		}

		postCollaborators[i] = &newPostCollaborator
	}

	return postCollaborators, nil
}
