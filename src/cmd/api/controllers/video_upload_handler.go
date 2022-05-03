package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/events"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"

	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	jsonDTO "github.com/Sogilis/Voogle/src/cmd/api/dto/json"
	protobufDTO "github.com/Sogilis/Voogle/src/cmd/api/dto/protobuf"
	"github.com/Sogilis/Voogle/src/cmd/api/metrics"
	"github.com/Sogilis/Voogle/src/cmd/api/models"
)

type VideoUploadHandler struct {
	S3Client      clients.IS3Client
	AmqpClient    clients.IAmqpClient
	MariadbClient *sql.DB
	UUIDGen       uuidgenerator.IUUIDGenerator
}

type Response struct {
	Video jsonDTO.VideoJson           `json:"video"`
	Links map[string]jsonDTO.LinkJson `json:"_links"`
}

// VideoUploadHandler godoc
// @Summary Upload video file
// @Description Upload video file
// @Tags upload
// @Accept multipart/form-data
// @Produce plain
// @Param file formData file true "video"
// @Success 200 {Json} Video and Links (HATEOAS)
// @Failure 400 {object} object
// @Failure 409 {object} object
// @Failure 415 {object} object
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
	sourceName := "source" + filepath.Ext(fileHandler.Filename)

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

	// Generate video UUID
	videoID, err := v.UUIDGen.GenerateUuid()
	if err != nil {
		log.Error("Cannot generate new videoID : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create new video
	videoCreated, err := dao.CreateVideo(r.Context(), v.MariadbClient, videoID, title, int(models.UPLOADING))
	if err != nil {
		// Check if the returned error comes from duplicate title
		videoCreated, err = dao.GetVideoFromTitle(r.Context(), v.MariadbClient, title)
		if err != nil {
			log.Error("Cannot find video ", title, "  : ", err)
			log.Error("Cannot insert new video into database: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info("This title already exist, check video status")
		if videoCreated.Status == models.FAIL_UPLOAD {
			// Retry to upload+encode
			log.Debugf("Last upload of video %v failed, simply retry", videoCreated.Title)

		} else if videoCreated.Status == models.FAIL_ENCODE {
			// Retry to encode
			log.Debug("Ask for video encoding")
			if err = v.sendVideoForEncoding(r.Context(), sourceName, videoCreated); err != nil {
				log.Error("Cannot send video for encoding : ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Here, everything went well: video encoding asked and video status updated
			return

		} else {
			// Title already exist, video already upload and encode, return error
			log.Error("A video with this title already uploaded and encoded")
			http.Error(w, "This title already exists", http.StatusConflict)
			return
		}
	}

	// Create new upload
	if err := v.uploadVideo(videoCreated, file, sourceName, r); err != nil {
		log.Error("Cannot upload vieo : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = v.sendVideoForEncoding(r.Context(), sourceName, videoCreated); err != nil {
		log.Error("Cannot send video for encoding : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Include videoCreated and HATEOAS upload link into response
	writeHTTPResponse(videoCreated, w)

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

func (v VideoUploadHandler) uploadVideo(video *models.Video, file multipart.File, sourceName string, r *http.Request) error {
	uploadID, err := v.UUIDGen.GenerateUuid()
	if err != nil {
		log.Error("Cannot generate new uploadID : ", err)
		return err
	}

	uploadCreated, err := dao.CreateUpload(r.Context(), v.MariadbClient, uploadID, video.ID, int(models.STARTED))
	if err != nil {
		// Update video status : FAIL_UPLOAD
		log.Error("Cannot insert new upload into database: ", err)
		video.Status = models.FAIL_UPLOAD
		if err := dao.UpdateVideo(r.Context(), v.MariadbClient, video); err != nil {
			log.Errorf("Unable to update video with status  %v: %v", video.Status, err)
		}
		return err
	}

	// Upload on S3
	err = v.S3Client.PutObjectInput(r.Context(), file, video.ID+"/"+sourceName)
	if err != nil {
		// Update video status : FAIL_UPLOAD + uploads status : FAILED
		log.Error("Unable to put object input on S3 ", err)

		if err := videoAndUploadFailed(r.Context(), video, uploadCreated, v.MariadbClient); err != nil {
			log.Error("video and upload status failed : ", err)
		}

		return err
	}
	log.Debug("Success upload video " + video.ID + " on S3")

	// Same time for videos and uploads
	uploadDate := time.Now()

	// Update videos status : UPLOADED + Upload date
	video.Status = models.UPLOADED
	video.UploadedAt = &uploadDate
	if err = dao.UpdateVideo(r.Context(), v.MariadbClient, video); err != nil {
		// Update video status : FAIL_UPLOAD + Update uploads status : FAILED
		log.Errorf("Unable to update video with status  %v : %v", video.Status, err)
		if err = videoAndUploadFailed(r.Context(), video, uploadCreated, v.MariadbClient); err != nil {
			log.Error("video and upload status failed : ", err)
			return err
		}

		err = v.S3Client.RemoveObject(r.Context(), video.ID+"/"+sourceName)
		if err != nil {
			log.Errorf("Unable to remove uploaded video  %v : %v", video.ID, err)
		}
		return err
	}

	// Update uploads status : DONE + Upload date
	uploadCreated.Status = models.DONE
	uploadCreated.UploadedAt = &uploadDate
	if err = dao.UpdateUpload(r.Context(), v.MariadbClient, uploadCreated); err != nil {
		// Update uploads status : FAILED
		log.Errorf("Unable to update upload with status  %v: %v", uploadCreated.Status, err)
		if err = videoAndUploadFailed(r.Context(), video, uploadCreated, v.MariadbClient); err != nil {
			log.Error("video and upload status failed : ", err)
			return err
		}

		err = v.S3Client.RemoveObject(r.Context(), video.ID+"/"+sourceName)
		if err != nil {
			log.Errorf("Unable to remove uploaded video  %v : %v", video.ID, err)
		}
		return err
	}

	return nil
}

func (v VideoUploadHandler) sendVideoForEncoding(ctx context.Context, sourceName string, video *models.Video) error {
	videoProto := protobufDTO.VideoToVideoProtobuf(video, sourceName)
	videoData, err := proto.Marshal(videoProto)
	if err != nil {
		log.Error("Unable to marshal video : ", err)
		videoEncodeFailed(ctx, video, v.MariadbClient)
		return err
	}

	if err = v.AmqpClient.Publish(events.VideoUploaded, videoData); err != nil {
		log.Error("Unable to publish on Amqp client : ", err)
		videoEncodeFailed(ctx, video, v.MariadbClient)
		return err
	}

	// Update video status : ENCODING
	video.Status = models.ENCODING
	if err := dao.UpdateVideo(ctx, v.MariadbClient, video); err != nil {
		log.Errorf("Unable to update video with status  %v: %v", video.Status, err)
		videoEncodeFailed(ctx, video, v.MariadbClient)
		return err
	}

	return nil
}

func videoEncodeFailed(ctx context.Context, video *models.Video, db *sql.DB) {
	// Update video status : FAIL_ENCODE
	video.Status = models.FAIL_ENCODE
	if err := dao.UpdateVideo(ctx, db, video); err != nil {
		log.Errorf("Unable to update video with status  %v: %v", video.Status, err)
	}
}

func videoAndUploadFailed(ctx context.Context, video *models.Video, upload *models.Upload, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Error("Cannot open new database transaction : ", err)
		return err
	}

	video.Status = models.FAIL_UPLOAD
	upload.Status = models.FAILED
	if err := dao.UpdateVideoTx(ctx, tx, video); err != nil {
		log.Errorf("Unable to update video with status  %v: %v", video.Status, err)
		if err := tx.Rollback(); err != nil {
			log.Error("Cannot rollback : ", err)
		}
		return err
	}
	if err := dao.UpdateUploadTx(ctx, tx, upload); err != nil {
		log.Errorf("Unable to update upload with status  %v: %v", upload.Status, err)
		if err := tx.Rollback(); err != nil {
			log.Error("Cannot rollback : ", err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error("Cannot commit database transaction")
		if err := tx.Rollback(); err != nil {
			log.Error("Cannot rollback : ", err)
		}
		return err
	}

	return nil
}

func writeHTTPResponse(video *models.Video, w http.ResponseWriter) {
	// Include videoCreated and status link into response (HATEOAS)

	videoJson := jsonDTO.VideoToVideoJson(video)

	links := map[string]jsonDTO.LinkJson{
		"status": jsonDTO.LinkToLinkJson(&models.Link{Href: "/api/v1/videos/" + video.ID + "/status", Method: "GET"}),
		"stream": jsonDTO.LinkToLinkJson(&models.Link{Href: "/api/v1/videos/" + video.ID + "/streams/master.m3u8", Method: "GET"}),
	}

	response := Response{
		Video: videoJson,
		Links: links,
	}

	payload, err := json.Marshal(response)
	if err != nil {
		log.Error("Unable to parse data struct in json ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(payload)

}
