package controllers

import (
	"database/sql"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	"github.com/Sogilis/Voogle/src/cmd/api/db/models"
	"github.com/Sogilis/Voogle/src/cmd/api/metrics"
	"github.com/Sogilis/Voogle/src/pkg/clients"
	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"
	"github.com/Sogilis/Voogle/src/pkg/events"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"
)

type VideoUploadHandler struct {
	S3Client      clients.IS3Client
	AmqpClient    clients.IAmqpClient
	MariadbClient *sql.DB
	UUIDGen       uuidgenerator.IUUIDGenerator
}

// VideoUploadHandler godoc
// @Summary Upload video file
// @Description Upload video file
// @Tags upload
// @Accept multipart/form-data
// @Produce plain
// @Param file formData file true "video"
// @Success 200 {string} string "OK"
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /api/v1/videos/upload [post]
func (v VideoUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	// Check if the received file is a supported video
	if typeOk := isSupportedType(file); !typeOk {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := v.UUIDGen.GenerateUuid()
	if err != nil {
		log.Error("Cannot generate new id : ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	publicId, err := v.UUIDGen.GenerateUuid()
	if err != nil {
		log.Error("Cannot generate new public id : ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Update database
	videoModel := models.VideoUpload{
		Title:    title,
		Id:       id,
		PublicId: publicId,
	}

	if err = dao.CreateVideo(v.MariadbClient, &videoModel); err != nil {
		log.Error("Cannot insert new video to database: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Upload on S3
	sourceName := "source" + filepath.Ext(fileHandler.Filename)
	err = v.S3Client.PutObjectInput(r.Context(), file, videoModel.PublicId+"/"+sourceName)
	if err != nil {
		log.Error("Unable to put object input on S3 ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Debug("Success upload video " + videoModel.PublicId + " on S3")

	video := &contracts.Video{
		Id:     videoModel.PublicId,
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
	metrics.CounterVideoUploadSuccess.Inc()
}

func isSupportedType(input io.ReaderAt) bool {
	// Use ReadAt instead of Read to avoid seek affect resulting
	// in readed bytes missing
	buff := make([]byte, 262) // 262 bytes : no need more for video format
	if nbByte, err := input.ReadAt(buff, 0); err != nil {
		log.Error("Cannot check file type : ", err)
		log.Error("Asked 262 bytes, readed : ", nbByte)
		return false
	}

	mime := mimetype.Detect(buff)
	log.Debug("Receive " + mime.Extension() + " file")

	if !strings.EqualFold(mime.String()[:5], "video") {
		log.Error("Not supported file type : " + mime.String() + " (" + mime.Extension() + ")")
		return false
	}
	return true
}
