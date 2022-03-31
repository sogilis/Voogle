package models

import "time"

type VideoStatus int

const (
	UPLOADING VideoStatus = iota
	UPLOADED
	ENCODING
	COMPLETE
	UNKNOWN
	FAIL_UPDLOAD
	FAIL_ENCODE
)

type Video struct {
	ID          string
	Title       string
	VideoStatus VideoStatus
	UploadedAt  *time.Time
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

type VideoUpload struct {
	Title    string
	Id       string
	PublicId string
}

type UploadStatus int

const (
	STARTED UploadStatus = iota
	DONE
	FAILED
)

type Upload struct {
	ID           string
	VideoId      string
	UploadStatus string
	UploadedAt   *time.Time
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
	Progress     int
}
