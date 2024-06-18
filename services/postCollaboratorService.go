package services

import (
	"fmt"
	"slices"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type PostCollaboratorService struct {
	PostCollaboratorRepository database.ModelRepositoryInterface[*models.PostCollaborator]
	MemberRepository           database.ModelRepositoryInterface[*models.Member]
	PostRepository             database.ModelRepositoryInterface[*models.Post]
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

// We add all branch collaborators to the project post as post collaborators with the "reviewer" type, unless they have already been added as such
func (postCollaboratorService *PostCollaboratorService) MergeReviewers(projectPost *models.ProjectPost, reviews []*models.BranchReview) error {
	// get all member ids which are reviewers present in post collaborators initially
	reviewerMemberIDs := []uint{}

	// load the post object (so we preload collaborators)
	post, err := postCollaboratorService.PostRepository.GetByID(projectPost.PostID)
	if err != nil {
		return fmt.Errorf("failed to get post with ID %d of project post with ID %d: %w", projectPost.PostID, projectPost.ID, err)
	}

	for _, c := range post.Collaborators {
		if c.CollaborationType == models.Reviewer {
			reviewerMemberIDs = append(reviewerMemberIDs, c.MemberID)
		}
	}

	// add all new post collaborators
	for _, review := range reviews {
		// if the member is already present as a post collaborator, we do not add it again
		if slices.Contains(reviewerMemberIDs, review.MemberID) {
			continue
		}

		// otherwise we add this post collaborator
		reviewMember, err := postCollaboratorService.MemberRepository.GetByID(review.MemberID)
		if err != nil {
			return fmt.Errorf("failed to get reviewing member from db")
		}

		asPostCollaborator := &models.PostCollaborator{
			Member:            *reviewMember,
			PostID:            projectPost.PostID,
			CollaborationType: models.Reviewer,
		}

		post.Collaborators = append(post.Collaborators, asPostCollaborator)
	}

	// Finally, re-assign the collaborators field
	projectPost.Post.Collaborators = post.Collaborators

	// Save changes to the post
	if _, err := postCollaboratorService.PostRepository.Update(post); err != nil {
		return fmt.Errorf("failed to update post metadata: %w", err)
	}

	return nil
}

// We add all branch collaborators to the project post as post collaborators with the "contributor" type, unless they have already been added as such
func (postCollaboratorService *PostCollaboratorService) MergeContributors(projectPost *models.ProjectPost, branchCollaborators []*models.BranchCollaborator) error {
	// get all member ids which are collaborators present in post collaborators initially
	contributorMemberIDs := []uint{}

	// load the post object (so we preload collaborators)
	post, err := postCollaboratorService.PostRepository.GetByID(projectPost.PostID)
	if err != nil {
		return fmt.Errorf("failed to get post with ID %d of project post with ID %d: %w", projectPost.PostID, projectPost.ID, err)
	}

	for _, c := range post.Collaborators {
		if c.CollaborationType == models.Contributor {
			contributorMemberIDs = append(contributorMemberIDs, c.MemberID)
		}
	}

	// add all new post collaborators
	for _, branchCollaborator := range branchCollaborators {
		// if the member is already present as a post collaborator, we do not add it again
		if slices.Contains(contributorMemberIDs, branchCollaborator.MemberID) {
			continue
		}

		// otherwise we add this post collaborator
		branchCollaboratorMember, err := postCollaboratorService.MemberRepository.GetByID(branchCollaborator.MemberID)
		if err != nil {
			return fmt.Errorf("failed to get contributing member from db")
		}

		asPostCollaborator := &models.PostCollaborator{
			Member:            *branchCollaboratorMember,
			PostID:            projectPost.PostID,
			CollaborationType: models.Contributor,
		}

		post.Collaborators = append(post.Collaborators, asPostCollaborator)
	}

	// Finally, re-assign the collaborators field
	projectPost.Post.Collaborators = post.Collaborators

	// Save changes to the post
	if _, err := postCollaboratorService.PostRepository.Update(post); err != nil {
		return fmt.Errorf("failed to update post metadata: %w", err)
	}

	return nil
}
