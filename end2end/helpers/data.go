package helpers

type VideoInfo struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}
type AllVideos struct {
	Status string      `json:"status"`
	Data   []VideoInfo `json:"data"`
}
