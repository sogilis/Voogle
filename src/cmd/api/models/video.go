package models

import (
	"errors"
	"strings"
	"time"
)

type VideoStatus int

const (
	UNSPECIFIED VideoStatus = iota
	UPLOADING
	UPLOADED
	ENCODING
	COMPLETE
	ARCHIVE
	UNKNOWN
	FAIL_UPLOAD
	FAIL_ENCODE
)

func (v VideoStatus) String() string {
	switch v {
	case UNSPECIFIED:
		return "Unspecified"
	case UPLOADING:
		return "Uploading"
	case UPLOADED:
		return "Uploaded"
	case ENCODING:
		return "Encoding"
	case COMPLETE:
		return "Complete"
	case ARCHIVE:
		return "Archive"
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

func StringToVideoStatus(v string) (VideoStatus, error) {
	switch strings.ToUpper(v) {
	case "UNSPECIFIED":
		return UNSPECIFIED, nil
	case "UPLOADING":
		return UPLOADING, nil
	case "UPLOADED":
		return UPLOADED, nil
	case "ENCODING":
		return ENCODING, nil
	case "COMPLETE":
		return COMPLETE, nil
	case "ARCHIVE":
		return ARCHIVE, nil
	case "UNKNOWN":
		return UNKNOWN, nil
	case "FAIL_UPLOAD":
		return FAIL_UPLOAD, nil
	case "FAIL_ENCODE":
		return FAIL_ENCODE, nil
	default:
		return UNSPECIFIED, errors.New("No cast for " + v + " to VideoStatus")
	}
}

type Video struct {
	ID         string
	Title      string
	Status     VideoStatus
	UploadedAt *time.Time
	CreatedAt  *time.Time
	UpdatedAt  *time.Time
	SourcePath string
	CoverPath  string
}
