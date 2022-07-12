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
	UploadedAt *time.Time `json:"uploadedAt" example:"2022-04-15T12:59:52Z"`
	CreatedAt  *time.Time `json:"createdAt" example:"2022-04-15T12:59:52Z"`
	UpdatedAt  *time.Time `json:"updatedAt" example:"2022-04-15T12:59:52Z"`
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
	Title  string `json:"title" example:"AmazingTitle"`
	Status string `json:"status" example:"UPLOADED"`
}

func VideoToStatusJson(video *models.Video) VideoStatus {
	videoStatus := VideoStatus{
		Title:  video.Title,
		Status: video.Status.String(),
	}

	return videoStatus
}

// VideoInfo DTO

type VideoInfo struct {
	Title          string `json:"title" example:"amazingtitle"`
	UploadDateUnix int64  `json:"uploadDateUnix" example:"1652173257"`
}

func VideoToInfoJson(video *models.Video) VideoInfo {
	videoInfo := VideoInfo{
		Title:          video.Title,
		UploadDateUnix: video.UploadedAt.Unix(),
	}

	return videoInfo
}

// LinkJson DTO

type LinkJson struct {
	Href   string `json:"href" exemple:"api/v1/videos/{id}/status"`
	Method string `json:"method" exemple:"GET"`
}

func LinkToLinkJson(link *models.Link) LinkJson {
	linkJson := LinkJson{
		Href:   link.Href,
		Method: link.Method,
	}

	return linkJson
}
