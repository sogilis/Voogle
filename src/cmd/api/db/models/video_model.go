package models

type VideoModel struct {
	Id         string
	ClientId   string
	Title      string
	VState     string
	LastUpdate string
}

type VideoModelUpload struct {
	Title string
}
