package models

import (
	"time"
)

type VideoStatus int

const (
	UNSPECIFIED VideoStatus = iota
	UPLOADING
	UPLOADED
	ENCODING
	COMPLETE
	UNKNOWN
	FAIL_UPLOAD
	FAIL_ENCODE
)

func (v VideoStatus) String() string {
	switch v {
	case UPLOADING:
		return "Uploading"
	case UPLOADED:
		return "Uploaded"
	case ENCODING:
		return "Encoding"
	case COMPLETE:
		return "Complete"
	case UNKNOWN:
		return "Unknown"
	case FAIL_UPLOAD:
		return "Fail_upload"
	case FAIL_ENCODE:
		return "Fail_encode"
	default:
		return "VideoStatus unspecified"
	}
}

type Video struct {
	ID         string
	Title      string
	Status     VideoStatus
	UploadedAt *time.Time
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
}
