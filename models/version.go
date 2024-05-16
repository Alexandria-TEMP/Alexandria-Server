package models

import "gorm.io/gorm"

type RenderStatus string

const (
	Success RenderStatus = "success"
	Pending RenderStatus = "pending"
	Failure RenderStatus = "failure"
)

type Version struct {
	gorm.Model

	// Version has many Discussion
	Discussions  []Discussion `gorm:"foreignKey:VersionID"`
	RenderStatus RenderStatus
}
