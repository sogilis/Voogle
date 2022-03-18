package models

type VideoModel struct {
	Id         []uint8
	ClientId   []uint8
	Title      string
	VState     string
	LastUpdate string
}

type VideoModelUpload struct {
	Title string
}
