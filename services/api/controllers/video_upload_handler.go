package controllers

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/services/api/clients"
)

type Video struct {
	Title string `json:"title"`
}

type VideoUploadHandler struct {
	S3Client    clients.IS3Client
	RedisClient clients.IRedisClient
}

func (v VideoUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("POST VideoUploadHandler")

	file, fileHandler, err := r.FormFile("video")
	if err != nil {
		log.Error("Missing file ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	if title == "" {
		log.Error("Missing title file ")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	title = strings.ReplaceAll(title, " ", "_")

	err = v.S3Client.PutObjectInput(r.Context(), file, title+"/source."+filepath.Ext(fileHandler.Filename))
	if err != nil {
		log.Error("Unable to put object input on S3 ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Debug("Success upload video " + title + " on S3")

	json, err := json.Marshal(Video{Title: title})
	if err != nil {
		log.Error("Unable to marshal video")
	}

	err = v.RedisClient.Publish(r.Context(), "video_uploaded_on_S3", json)
	if err != nil {
		log.Error("Unable to publish on Redis client")
	}
}
