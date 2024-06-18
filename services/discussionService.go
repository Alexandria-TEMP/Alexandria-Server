package services

import (
	"fmt"

	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/database"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/forms"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

type DiscussionService struct {
	DiscussionRepository          database.ModelRepositoryInterface[*models.Discussion]
	DiscussionContainerRepository database.ModelRepositoryInterface[*models.DiscussionContainer]
	MemberRepository              database.ModelRepositoryInterface[*models.Member]
}

func (discussionService *DiscussionService) GetDiscussion(id uint) (*models.Discussion, error) {
	return discussionService.DiscussionRepository.GetByID(id)
}

func (discussionService *DiscussionService) CreateRootDiscussion(form *forms.RootDiscussionCreationForm, member *models.Member) (*models.Discussion, error) {
	// Verify the target container exists
	_, err := discussionService.DiscussionContainerRepository.GetByID(form.ContainerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get parent discussion container: %w", err)
	}

	// If the discussion is anonymous, ignore the member ID field, and leave 'member' blank
	anonymous := form.DiscussionCreationForm.Anonymous
	if anonymous {
		member = nil
	}

	// The discussion to be created
	discussion := models.Discussion{
		ContainerID: form.ContainerID,
		Member:      member,
		Replies:     []*models.Discussion{},
		Text:        form.DiscussionCreationForm.Text,
	}

	// Try to create the discussion in the database
	if err := discussionService.DiscussionRepository.Create(&discussion); err != nil {
		return nil, err
	}

	return &discussion, nil
}

func (discussionService *DiscussionService) CreateReply(form *forms.ReplyDiscussionCreationForm, member *models.Member) (*models.Discussion, error) {
	// Verify the target parent discussion exists
	parentDiscussion, err := discussionService.DiscussionRepository.GetByID(form.ParentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get parent discussion: %w", err)
	}

	// If the discussion is anonymous, ignore the member ID field, and leave 'member' blank
	anonymous := form.DiscussionCreationForm.Anonymous
	if anonymous {
		member = nil
	}

	discussion := models.Discussion{
		ContainerID: parentDiscussion.ContainerID,
		Member:      member,
		Replies:     []*models.Discussion{},
		ParentID:    &form.ParentID,
		Text:        form.DiscussionCreationForm.Text,
	}

	// Try to create the discussion in the database
	if err := discussionService.DiscussionRepository.Create(&discussion); err != nil {
		return nil, err
	}

	return &discussion, nil
}
