package interfaces

import (
	"github.com/gofrs/flock"
	"gitlab.ewi.tudelft.nl/cse2000-software-project/2023-2024/cluster-v/17b/alexandria-backend/models"
)

//go:generate mockgen -package=mocks -source=./renderService_interface.go -destination=../../mocks/renderService_mock.go

type RenderService interface {
	// GetRender returns filepath of rendered repository ON A NON-MASTER BRANCH.
	// Error 1 is for status 202.
	// Error 2 is for status 404.
	GetRenderFile(branchID uint) (string, *flock.Flock, error, error)

	// GetRender returns filepath of rendered repository ON MAIN.
	// Error 1 is for status 202.
	// Error 2 is for status 404.
	GetMainRenderFile(postID uint) (string, *flock.Flock, error, error)

	// RenderBranch first unzips a zipped project, validates it, configures it to render well, renders it, and checks it rendered well.
	// We are handed a locked lock from the previous step in the render pipeline. It must be unlocked after rendering.
	RenderBranch(*models.Branch, *flock.Flock)

	// RenderPost first unzips a zipped project,validates it, configures it to render well, renders it, and checks it rendered well.
	// We are handed a locked lock from the previous step in the render pipeline. It must be unlocked after rendering.
	RenderPost(*models.Post, *flock.Flock)
}
