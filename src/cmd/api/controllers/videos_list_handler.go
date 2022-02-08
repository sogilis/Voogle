package controllers

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/cmd/api/clients"
)

type VideoInfo struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}
type AllVideos struct {
	Status string      `json:"status"`
	Data   []VideoInfo `json:"data"`
}

type VideosListHandler struct {
	S3Client clients.IS3Client
}

func (v VideosListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("GET VideosListHandler")

	videos, err := v.S3Client.ListObjects(r.Context())
	if err != nil {
		log.Error("Unable to list objects on S3", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	allVideos := AllVideos{}
	for _, video := range videos {
		videoInfo := VideoInfo{
			video,
			video,
		}
		allVideos.Data = append(allVideos.Data, videoInfo)
	}
	allVideos.Status = "Success"

	payload, err := json.Marshal(allVideos)

	if err != nil {
		log.Error("Unable to parse data struct in json", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(payload)
}
