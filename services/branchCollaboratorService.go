package services

import (
	"fmt"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type BranchCollaboratorService struct {
	BranchCollaboratorRepository database.ModelRepositoryInterface[*models.BranchCollaborator]
	MemberRepository             database.ModelRepositoryInterface[*models.Member]
}

func (branchCollaboratorService *BranchCollaboratorService) GetBranchCollaborator(id uint) (*models.BranchCollaborator, error) {
	return branchCollaboratorService.BranchCollaboratorRepository.GetByID(id)
}

func (branchCollaboratorService *BranchCollaboratorService) MembersToBranchCollaborators(memberIDs []uint, anonymous bool) ([]*models.BranchCollaborator, error) {
	// If anonymous, immediately return an empty list
	if anonymous {
		return []*models.BranchCollaborator{}, nil
	}

	// If not anonymous, ensure there's >= 1 member
	if len(memberIDs) < 1 {
		return nil, fmt.Errorf("failed to create branch collaborators: need at least 1 member ID")
	}

	// Create the branch collaborators
	branchCollaborators := make([]*models.BranchCollaborator, len(memberIDs))

	for index, memberID := range memberIDs {
		member, err := branchCollaboratorService.MemberRepository.GetByID(memberID)
		if err != nil {
			return nil, fmt.Errorf("failed to create branch collaborators: %w", err)
		}

		branchCollaborators[index] = &models.BranchCollaborator{
			Member: *member,
		}
	}

	return branchCollaborators, nil
}
