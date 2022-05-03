package json

import (
	"time"

	"github.com/Sogilis/Voogle/src/cmd/api/models"
)

// VideoJson DTO
type VideoJson struct {
	ID         string     `json:"id" example:"aaaa-b56b-..."`
	Title      string     `json:"title" example:"A Title"`
	Status     string     `json:"status" example:"VIDEO_STATUS_ENCODING"`
	UploadedAt *time.Time `json:"uploadedat" example:"2022-04-15T12:59:52Z"`
	CreatedAt  *time.Time `json:"createdat" example:"2022-04-15T12:59:52Z"`
	UpdatedAt  *time.Time `json:"updatedat" example:"2022-04-15T12:59:52Z"`
}

func VideoToVideoJson(video *models.Video) VideoJson {
	videoJson := VideoJson{
		ID:         video.ID,
		Title:      video.Title,
		Status:     video.Status.String(),
		CreatedAt:  video.CreatedAt,
		UploadedAt: video.UploadedAt,
		UpdatedAt:  video.UpdatedAt,
	}

	return videoJson
}

// VideoStatus DTO
type VideoStatus struct {
	Status string `json:"status" example:"UPLOADED"`
}

func VideoToStatusJson(video *models.Video) VideoStatus {
	videoStatus := VideoStatus{
		Status: video.Status.String(),
	}

	return videoStatus
}

// VideoInfo DTO

type VideoInfo struct {
	Title       string `json:"title" example:"amazingtitle"`
	Upload_date string `json:"uploaddate" example:"amazingtitle"`
}

func VideoToInfoJson(video *models.Video) VideoInfo {
	videoInfo := VideoInfo{
		Title:       video.Title,
		Upload_date: video.UploadedAt.Format("January 2, 2006 15:04:05"),
	}

	return videoInfo
}
