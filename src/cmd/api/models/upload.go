package models

import "time"

type UploadStatus int

const (
	STARTED UploadStatus = iota
	DONE
	FAILED
)

type Upload struct {
	ID         string
	VideoId    string
	Status     UploadStatus
	UploadedAt *time.Time
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	Progress   int
}
