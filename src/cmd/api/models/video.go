package models

import "time"

type VideoStatus int

const (
	UPLOADING VideoStatus = iota
	UPLOADED
	ENCODING
	COMPLETE
	UNKNOWN
	FAIL_UPLOAD
	FAIL_ENCODE
)

type Video struct {
	ID         string
	Title      string
	Status     VideoStatus
	UploadedAt *time.Time
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}

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
