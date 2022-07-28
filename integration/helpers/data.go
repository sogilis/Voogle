package helpers

import "time"

type VideoListResponse struct {
	Videos   []VideoInfo `json:"videos"`
	LastPage int         `json:"_lastpage"`
}

type VideoJson struct {
	ID         string     `json:"id" example:"aaaa-b56b-..."`
	Title      string     `json:"title" example:"A Title"`
	Status     string     `json:"status" example:"VIDEO_STATUS_ENCODING"`
	UploadedAt *time.Time `json:"uploadedAt" example:"2022-04-15T12:59:52Z"`
	CreatedAt  *time.Time `json:"createdAt" example:"2022-04-15T12:59:52Z"`
	UpdatedAt  *time.Time `json:"updatedAt" example:"2022-04-15T12:59:52Z"`
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
	Title  string `json:"title" example:"AmazingTitle"`
	Status string `json:"status" example:"UPLOADED"`
}

type VideoInfo struct {
	Title          string `json:"title" example:"amazingtitle"`
	UploadDateUnix int64  `json:"uploadDateUnix" example:"1652173257"`
}

type TransformerServiceListResponse struct {
	Services []TransformerServiceJson `json:"services"`
}

type TransformerServiceJson struct {
	Name string `json:"name" example:"gray"`
}
