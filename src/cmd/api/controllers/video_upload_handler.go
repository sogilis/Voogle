package controllers

import (
	"database/sql"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"
	"github.com/Sogilis/Voogle/src/pkg/events"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"

	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	"github.com/Sogilis/Voogle/src/cmd/api/metrics"
	"github.com/Sogilis/Voogle/src/cmd/api/models"
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

	// Check if the received file is a supported video type
	if typeOk := isSupportedType(file); !typeOk {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// Generate video and upload UUID
	videoID, err := v.UUIDGen.GenerateUuid()
	if err != nil {
		log.Error("Cannot generate new videoID : ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	uploadID, err := v.UUIDGen.GenerateUuid()
	if err != nil {
		log.Error("Cannot generate new uploadID : ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create new video
	videoCreated, err := dao.CreateVideo(v.MariadbClient, videoID, title, int(contracts.Video_UPLOADING))
	if err != nil {
		// Check if the returned error comes from duplicate title
		videoCreated, err = dao.GetVideoFromTitle(v.MariadbClient, title)
		if err != nil {
			log.Error("Cannot find video ", title, "  : ", err)
			log.Error("Cannot insert new video into database: ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Info("This title already exist, check video status")
		if videoCreated.Status == contracts.Video_FAIL_UPLOAD {
			// Retry to upload+encode
			log.Debugf("Last upload of video %v failed, simply retry", videoCreated.Title)

		} else if videoCreated.Status == contracts.Video_FAIL_ENCODE {
			// Retry to encode
			log.Debug("Ask for video encoding")
			video := &contracts.Video{
				Id:     videoCreated.ID,
				Source: "source" + filepath.Ext(fileHandler.Filename),
			}

			if err = sendVideoForEncoding(video, v.AmqpClient, videoCreated, v.MariadbClient); err != nil {
				log.Error("Cannot send video for encoding : ", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Here, everything went well: video encoding asked and video status updated
			return

		} else {
			// Title already exist, video already upload and encode, return error
			log.Error("A video with this title already uploaded and encoded")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// Create new upload
	uploadCreated, err := dao.CreateUpload(v.MariadbClient, uploadID, videoID, int(models.STARTED))
	if err != nil {
		log.Error("Cannot insert new upload into database: ", err)

		// Update video status : FAIL_UPLOAD
		videoUploadFailed(videoCreated, v.MariadbClient)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Upload on S3
	sourceName := "source" + filepath.Ext(fileHandler.Filename)
	err = v.S3Client.PutObjectInput(r.Context(), file, videoCreated.ID+"/"+sourceName)
	if err != nil {
		log.Error("Unable to put object input on S3 ", err)

		// Update video status : FAIL_UPLOAD
		videoUploadFailed(videoCreated, v.MariadbClient)

		// Update uploads status : FAILED
		uploadFailed(uploadCreated, v.MariadbClient)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Debug("Success upload video " + videoCreated.ID + " on S3")

	// Same time for videos and uploads
	uploadDate := time.Now()

	// Update videos status : UPLOADED + Upload date
	videoCreated.Status = contracts.Video_UPLOADED
	videoCreated.UploadedAt = &uploadDate
	if err = dao.UpdateVideo(v.MariadbClient, videoCreated); err != nil {
		log.Errorf("Unable to update video with status  %v : %v", videoCreated.Status, err)

		// Update video status : FAIL_UPLOAD
		videoUploadFailed(videoCreated, v.MariadbClient)

		// Update uploads status : FAILED
		uploadFailed(uploadCreated, v.MariadbClient)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update uploads status : DONE + Upload date
	uploadCreated.Status = models.DONE
	uploadCreated.UploadedAt = &uploadDate
	if err = dao.UpdateUpload(v.MariadbClient, uploadCreated); err != nil {
		log.Errorf("Unable to update upload with status  %v: %v", uploadCreated.Status, err)

		// Update uploads status : FAILED
		uploadFailed(uploadCreated, v.MariadbClient)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	video := &contracts.UploadedVideo{
		Id:     videoCreated.ID,
		Source: sourceName,
	}

	if err = sendVideoForEncoding(video, v.AmqpClient, videoCreated, v.MariadbClient); err != nil {
		log.Error("Cannot send video for encoding : ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//TODO : Include videoCreated into response
	//TODO : Include HATEOAS upload link
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

func videoUploadFailed(videoCreated *models.Video, db *sql.DB) {
	// Update video status : FAIL_UPLOAD
	videoCreated.Status = contracts.Video_FAIL_UPLOAD
	if err := dao.UpdateVideo(db, videoCreated); err != nil {
		log.Errorf("Unable to update video with status  %v: %v", videoCreated.Status, err)
	}
}

func uploadFailed(uploadCreated *models.Upload, db *sql.DB) {
	// Update upload status : FAILED
	uploadCreated.Status = models.FAILED
	if err := dao.UpdateUpload(db, uploadCreated); err != nil {
		log.Errorf("Unable to update upload with status  %v: %v", uploadCreated.Status, err)
	}
}

func sendVideoForEncoding(video *contracts.Video, amqpC clients.IAmqpClient, videoCreated *models.Video, db *sql.DB) error {
	videoData, err := proto.Marshal(video)
	if err != nil {
		log.Error("Unable to marshal video : ", err)
		return err
	}

	if err = amqpC.Publish(events.VideoUploaded, videoData); err != nil {
		log.Error("Unable to publish on Amqp client : ", err)
		return err
	}

	// Update video status : ENCODING
	videoCreated.Status = contracts.Video_ENCODING
	if err := dao.UpdateVideo(db, videoCreated); err != nil {
		log.Errorf("Unable to update video with status  %v: %v", videoCreated.Status, err)
		return err
	}

	metrics.CounterVideoUploadSuccess.Inc()
	return nil
}
