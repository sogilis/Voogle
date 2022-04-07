package models

import (
	"time"

	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"
)

type Video struct {
	ID         string
	Title      string
	Status     contracts.Video_VideoStatus
	UploadedAt *time.Time
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}
