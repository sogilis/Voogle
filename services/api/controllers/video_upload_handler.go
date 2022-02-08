package controllers

import (
	"net/http"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/Sogilis/Voogle/services/api/clients"
	contracts "github.com/Sogilis/Voogle/services/api/contracts/v1"
)

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

	sourceName := "source" + filepath.Ext(fileHandler.Filename)
	err = v.S3Client.PutObjectInput(r.Context(), file, title+"/"+sourceName)
	if err != nil {
		log.Error("Unable to put object input on S3 ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Debug("Success upload video " + title + " on S3")

	video := &contracts.Video{
		Id:     title,
		Source: sourceName,
	}

	videoData, err := proto.Marshal(video)
	if err != nil {
		log.Error("Unable to marshal video")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := v.RedisClient.Publish(r.Context(), "video_uploaded_on_S3", videoData); err != nil {
		log.Error("Unable to publish on Redis client")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
