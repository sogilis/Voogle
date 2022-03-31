package models

type Video struct {
	Id         string
	PublicId   string
	Title      string
	VState     string
	LastUpdate string
}

type VideoUpload struct {
	Title    string
	Id       string
	PublicId string
}
