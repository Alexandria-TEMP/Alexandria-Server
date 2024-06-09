package interfaces

import "gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"

//go:generate mockgen -package=mocks -source=./branchCollaboratorService_interface.go -destination=../../mocks/branchCollaboratorService_mock.go

type BranchCollaboratorService interface {
	GetBranchCollaborator(id uint) (*models.BranchCollaborator, error)
	MembersToBranchCollaborators(memberIDs []uint, anonymous bool) ([]*models.BranchCollaborator, error)
}
