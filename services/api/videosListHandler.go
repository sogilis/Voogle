package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type videoInfo struct {
	Title string `json:"title"`
}
type allVideos struct {
	Status string      `json:"status"`
	Data   []videoInfo `json:"data"`
}

type videosListHandler struct{}

func (v videosListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = os.Mkdir("./videos", os.ModePerm)
	files, err := ioutil.ReadDir("./videos")
	if err != nil {
		log.Error("Unable to resolve directory videos", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	allVideos := allVideos{}

	for _, f := range files {
		videoInfo := videoInfo{
			f.Name(),
		}
		allVideos.Data = append(allVideos.Data, videoInfo)
	}
	allVideos.Status = "Success"

	payload, err := json.Marshal(allVideos)

	if err != nil {
		log.Error("Unable to parse data struc in json", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(payload)
}
