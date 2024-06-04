package interfaces

type RenderService interface {
	// GetRender returns filepath of rendered repository.
	// Error 1 is for status 202.
	// Error 2 is for status 404.
	GetRenderFile(branchID uint) (string, error, error)

	// Render first unzips a zipper project, validates it, renders it, and checks it rendered well.
	Render()
}
