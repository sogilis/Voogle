package helpers

import "time"

type VideoInfo struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}
type AllVideos struct {
	Status string      `json:"status"`
	Data   []VideoInfo `json:"data"`
}

type VideoJson struct {
	ID         string     `json:"id" example:"aaaa-b56b-..."`
	Title      string     `json:"title" example:"A Title"`
	Status     string     `json:"status" example:"VIDEO_STATUS_ENCODING"`
	UploadedAt *time.Time `json:"uploadedat" example:"2022-04-15T12:59:52Z"`
	CreatedAt  *time.Time `json:"createdat" example:"2022-04-15T12:59:52Z"`
	UpdatedAt  *time.Time `json:"updatedat" example:"2022-04-15T12:59:52Z"`
}

type Link struct {
	Rel    string `json:"rel" example:"getStatus"`
	Href   string `json:"href" example:"api/v0/..."`
	Method string `json:"method" example:"GET"`
}

type Response struct {
	Video VideoJson `json:"video"`
	Links []Link    `json:"links"`
}

type VideoStatus struct {
	Status string `json:"status" example:"UPLOADED"`
}
