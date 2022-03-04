package controllers

import (
	"errors"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	metrics "github.com/Sogilis/Voogle/src/cmd/api/metrics"
	"github.com/Sogilis/Voogle/src/pkg/clients"
	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"
	"github.com/Sogilis/Voogle/src/pkg/events"
)

type VideoUploadHandler struct {
	S3Client   clients.IS3Client
	AmqpClient clients.IAmqpClient
}

func (v VideoUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	metrics.CounterApiVideoUploadInit.Inc()
	log.Debug("POST VideoUploadHandler")

	file, fileHandler, err := r.FormFile("video")
	if err != nil {
		log.Error("Missing file ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	title := r.FormValue("title")
	if title == "" {
		log.Error("Missing title file ")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	title = strings.ReplaceAll(title, " ", "_")

	// Check if the received file is a supported video
	if err := isSupportedType(&file); err != nil {
		log.Error("File isn't a video : ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	if err = v.AmqpClient.Publish(events.VideoUploaded, videoData); err != nil {
		log.Error("Unable to publish on Amqp client ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	metrics.CounterApiVideoUploadSuccess.Inc()
}

func isSupportedType(filep *multipart.File) error {
	file := *filep
	buff := make([]byte, 262) // 262 bytes : no need more for video format
	if _, err := file.Read(buff); err != nil {
		return err
	}

	mtype := mimetype.Detect(buff)
	log.Debug("Receive " + mtype.Extension() + " file")

	//FFMPEG doesn't support video/quicktime (mov, mqv)
	if !strings.EqualFold(mtype.String()[:5], "video") || mtype.Is("video/quicktime") {
		return errors.New("wrong file type")
	}
	return nil
}
