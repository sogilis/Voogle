package controllers

import (
	"net/http"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"
	"github.com/Sogilis/Voogle/src/pkg/events"
)

type VideoUploadHandler struct {
	S3Client       clients.IS3Client
	RabbitmqClient clients.IRabbitmqClient
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

	if err = v.RabbitmqClient.Publish(events.VideoUploaded, videoData); err != nil {
		// if err = v.RabbitmqClient.Connect(); err != nil {
		// 	log.Error("Unable to publish on Rabbitmq client ", err)
		// 	return
		// }

		// if _, err = v.RabbitmqClient.QueueDeclare(events.VideoUploaded); err != nil {
		// 	log.Error("Unable to publish on Rabbitmq client ", err)
		// 	return
		// }

		// if err = v.RabbitmqClient.Publish(events.VideoUploaded, videoData); err != nil {
		log.Error("Unable to publish on Rabbitmq client ", err)
		w.WriteHeader(http.StatusInternalServerError)
		// }
		return
	}
}
