package models

type VideoModel struct {
	Id         string
	PublicId   string
	Title      string
	VState     string
	LastUpdate string
}

type VideoModelUpload struct {
	Title    string
	Id       string
	PublicId string
}
