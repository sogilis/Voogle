package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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

	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	jsonDTO "github.com/Sogilis/Voogle/src/cmd/api/dto/json"
	"github.com/Sogilis/Voogle/src/cmd/api/dto/protobuf"
	protobufDTO "github.com/Sogilis/Voogle/src/cmd/api/dto/protobuf"
	"github.com/Sogilis/Voogle/src/cmd/api/metrics"
	"github.com/Sogilis/Voogle/src/cmd/api/models"
)

type VideoUploadHandler struct {
	S3Client            clients.IS3Client
	AmqpClient          clients.IAmqpClient
	AmqpExchangerStatus clients.IAmqpExchanger
	VideosDAO           *dao.VideosDAO
	UploadsDAO          *dao.UploadsDAO
	UUIDGen             clients.IUUIDGenerator
}

type Response struct {
	Video jsonDTO.VideoJson           `json:"video"`
	Links map[string]jsonDTO.LinkJson `json:"_links"`
}

// VideoUploadHandler godoc
// @Summary Upload video file
// @Description Upload video file
// @Tags video
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "video"
// @Success 200 {object} Response "Video and Links (HATEOAS)"
// @Failure 400 {string} string
// @Failure 409 {string} string "This title already exists"
// @Failure 415 {string} string
// @Failure 500 {string} string
// @Router /api/v1/videos/upload [post]
func (v VideoUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) { //nolint:cyclop
	log.Debug("POST VideoUploadHandler")

	// Fetch title
	title := r.FormValue("title")
	if title == "" {
		log.Error("Missing title file ")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Infof("Receive video upload request with title : '%v'", title)

	// Fetch video
	fileVideo, fileHandler, err := r.FormFile("video")
	if err != nil {
		log.Error("Missing file ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer fileVideo.Close()

	// Check if the received file is a supported video type
	if !isSupportedVideoType(fileVideo) {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// Fetch cover image. Not mandatory
	fileCover, fileHandlerCover, err := r.FormFile("cover")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		log.Error("File cover error ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if fileCover != nil {
		defer fileCover.Close()

		// Check if the received file cover is a supported image type
		if !isSupportedCoverType(fileCover) {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
	}

	// Check if a video with this title already exists
	video, err := v.VideosDAO.GetVideoFromTitle(r.Context(), title)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if video != nil {
		// If a video with the same title already exists, and if its status is failed upload/encode,
		// try to re-upload/re-encode as needed
		if video.Status == models.FAIL_UPLOAD || video.Status == models.FAIL_ENCODE {
			v.resumeVideoUpload(r.Context(), video, fileCover, fileVideo, fileHandlerCover, w)
			return
		} else {
			// Title already exist, video already uploaded and encoded, return error
			log.Error("A video with this title already uploaded and encoded")
			http.Error(w, "This title already exists", http.StatusConflict)
			return
		}
	}

	// Generate video UUID
	videoID, err := v.UUIDGen.GenerateUuid()
	if err != nil {
		log.Error("Cannot generate new video ID : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Upload cover image (if exists) on S3, update database
	coverPath, err := v.uploadCover(r.Context(), fileCover, videoID, fileHandlerCover)
	if err != nil {
		log.Error("Cannot upload cover image : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Upload video on S3, update database
	videoPath := videoID + "/" + "source" + filepath.Ext(fileHandler.Filename)
	videoCreated, err := v.uploadVideo(r.Context(), videoID, title, videoPath, coverPath, fileVideo, nil)
	if err != nil {
		log.Error("Cannot upload video : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = v.sendVideoForEncoding(r.Context(), videoCreated); err != nil {
		log.Error("Cannot send video for encoding : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Include videoCreated and HATEOAS upload link into response
	writeHTTPResponse(videoCreated, w)
	log.Infof("Video '%v' successfully uploaded", title)
}

func isSupportedVideoType(input io.ReaderAt) bool {
	// Use ReadAt instead of Read to avoid seek affect resulting in readed bytes missing
	buff := make([]byte, 262) // 262 bytes : need no more for video format
	if _, err := input.ReadAt(buff, 0); err != nil {
		log.Error("Cannot check file type : ", err)
		return false
	}

	mime := mimetype.Detect(buff)
	if !strings.EqualFold(mime.String()[:5], "video") {
		log.Error("Unsupported file type : " + mime.String() + " (" + mime.Extension() + ")")
		return false
	}
	return true
}

func isSupportedCoverType(input io.ReaderAt) bool {
	// Use ReadAt instead of Read to avoid seek affect resulting in readed bytes missing
	buff := make([]byte, 262) // 262 bytes : no need more for image format
	if _, err := input.ReadAt(buff, 0); err != nil {
		log.Error("Cannot check file type : ", err)
		return false
	}

	mime := mimetype.Detect(buff)
	if !strings.EqualFold(mime.String()[:5], "image") {
		log.Error("Not supported file type : " + mime.String() + " (" + mime.Extension() + ")")
		return false
	}

	for _, ext := range []string{".jpeg", ".jpg", ".png"} {
		if strings.EqualFold(mime.Extension(), ext) {
			return true
		}
	}

	log.Error("Not supported file type : " + mime.String() + " (" + mime.Extension() + ")")
	return false
}

func (v VideoUploadHandler) resumeVideoUpload(ctx context.Context, video *models.Video, fileCover, fileVideo multipart.File, fileHandler *multipart.FileHeader, w http.ResponseWriter) {

	// If the upload failed before the encoding started, then we have to fix the upload before resuming with the encoding.
	if video.Status == models.FAIL_UPLOAD {
		log.Debug("Try to re-upload failed video")
		coverPath, err := v.uploadCover(ctx, fileCover, video.ID, fileHandler)
		if err != nil {
			log.Error("Cannot upload cover image : ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		video, err = v.uploadVideo(ctx, video.ID, video.Title, video.SourcePath, coverPath, fileVideo, video)
		if err != nil {
			log.Error("Cannot upload video : ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	log.Debug("Try to re-encode failed video")
	if err := v.sendVideoForEncoding(ctx, video); err != nil {
		log.Error("Cannot send video for encoding : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Include video and HATEOAS upload link into response
	writeHTTPResponse(video, w)
	log.Infof("Video '%v' successfully uploaded", video.Title)
}

func (v VideoUploadHandler) uploadCover(ctx context.Context, cover multipart.File, videoID string, fileHandler *multipart.FileHeader) (string, error) {
	coverPath := ""
	if cover != nil {
		coverPath = videoID + "/" + "cover" + filepath.Ext(fileHandler.Filename)
		if err := v.S3Client.PutObjectInput(ctx, cover, coverPath); err != nil {
			log.Error("Cannot upload cover : ", err)
			return "", err
		}
	}
	return coverPath, nil
}

func (v VideoUploadHandler) uploadVideo(ctx context.Context, videoID, title, videoPath, coverPath string, file multipart.File, video *models.Video) (*models.Video, error) {
	metrics.CounterVideoUploadRequest.Inc()

	// video not nil means that the video already exists. So we are in case of recover after error
	if video == nil {
		var err error
		video, err = v.VideosDAO.CreateVideo(ctx, videoID, title, int(models.UPLOADING), videoPath, coverPath)
		if err != nil {
			metrics.CounterVideoUploadFail.Inc()
			log.Error("Cannot generate new uploadID : ", err)

			return nil, err
		}
	}

	uploadID, err := v.UUIDGen.GenerateUuid()
	if err != nil {
		metrics.CounterVideoUploadFail.Inc()
		log.Error("Cannot generate new uploadID : ", err)

		return nil, err
	}

	uploadCreated, err := v.UploadsDAO.CreateUpload(ctx, uploadID, video.ID, int(models.STARTED))
	if err != nil {
		metrics.CounterVideoUploadFail.Inc()
		log.Error("Cannot insert new upload into database: ", err)

		v.videoUploadFailed(ctx, video)
		v.publishStatus(video)
		return nil, err
	}

	// Upload video on S3
	err = v.S3Client.PutObjectInput(ctx, file, video.SourcePath)
	if err != nil {
		metrics.CounterVideoUploadFail.Inc()
		log.Error("Unable to put object input on S3 ", err)

		if err := v.videoAndUploadFailed(ctx, video, uploadCreated); err != nil {
			log.Error("video and upload status failed : ", err)
			return nil, err
		}

		return nil, err
	}
	log.Debug("Success upload video " + video.ID + " on S3")

	// Same time for videos and uploads
	uploadDate := time.Now()

	// Update videos status : UPLOADED + Upload date
	video.Status = models.UPLOADED
	video.UploadedAt = &uploadDate
	if err = v.VideosDAO.UpdateVideo(ctx, video); err != nil {
		metrics.CounterVideoUploadFail.Inc()
		log.Errorf("Unable to update video with status  %v : %v", video.Status, err)

		if err := v.videoAndUploadFailed(ctx, video, uploadCreated); err != nil {
			log.Error("video and upload status failed : ", err)
			return nil, err
		}

		if err := v.S3Client.RemoveObject(ctx, video.SourcePath); err != nil {
			log.Errorf("Unable to remove uploaded video  %v : %v", video.ID, err)
			return nil, err
		}

		return nil, err
	}

	v.publishStatus(video)

	// Update uploads status : DONE + Upload date
	uploadCreated.Status = models.DONE
	uploadCreated.UploadedAt = &uploadDate
	if err = v.UploadsDAO.UpdateUpload(ctx, uploadCreated); err != nil {
		metrics.CounterVideoUploadFail.Inc()
		log.Errorf("Unable to update upload with status  %v: %v", uploadCreated.Status, err)

		if err := v.videoAndUploadFailed(ctx, video, uploadCreated); err != nil {
			log.Error("video and upload status failed : ", err)
			return nil, err
		}

		if err := v.S3Client.RemoveObject(ctx, video.SourcePath); err != nil {
			log.Errorf("Unable to remove uploaded video  %v : %v", video.ID, err)
			return nil, err
		}

		return nil, err
	}

	metrics.CounterVideoUploadSuccess.Inc()
	return video, nil
}

func (v VideoUploadHandler) sendVideoForEncoding(ctx context.Context, video *models.Video) error {
	metrics.CounterVideoEncodeRequest.Inc()

	videoProto := protobufDTO.VideoToVideoProtobuf(video)
	videoData, err := proto.Marshal(videoProto)
	if err != nil {
		metrics.CounterVideoEncodeFail.Inc()
		log.Error("Unable to marshal video : ", err)

		v.videoEncodeFailed(ctx, video)
		return err
	}

	if err := v.AmqpClient.Publish(events.VideoUploaded, videoData); err != nil {
		metrics.CounterVideoEncodeFail.Inc()
		log.Error("Unable to publish on Amqp client : ", err)

		v.videoEncodeFailed(ctx, video)
		return err
	}

	// Update video status : ENCODING
	video.Status = models.ENCODING
	if err := v.VideosDAO.UpdateVideo(ctx, video); err != nil {
		metrics.CounterVideoEncodeFail.Inc()
		log.Errorf("Unable to update video with status  %v: %v", video.Status, err)

		v.videoEncodeFailed(ctx, video)
		return err
	}

	v.publishStatus(video)

	return nil
}

func (v VideoUploadHandler) videoUploadFailed(ctx context.Context, video *models.Video) {
	video.Status = models.FAIL_UPLOAD
	if err := v.VideosDAO.UpdateVideo(ctx, video); err != nil {
		log.Errorf("Unable to update video with status  %v: %v", video.Status, err)
	}
}
func (v VideoUploadHandler) videoEncodeFailed(ctx context.Context, video *models.Video) {
	// Update video status : FAIL_ENCODE
	video.Status = models.FAIL_ENCODE
	if err := v.VideosDAO.UpdateVideo(ctx, video); err != nil {
		log.Errorf("Unable to update video with status  %v: %v", video.Status, err)
	}
}

func (v VideoUploadHandler) videoAndUploadFailed(ctx context.Context, video *models.Video, upload *models.Upload) error {
	tx, err := v.VideosDAO.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Error("Cannot open new database transaction : ", err)
		return err
	}

	// Defer a rollback in case anything fails.
	defer func() {
		_ = tx.Rollback()
	}()

	video.Status = models.FAIL_UPLOAD
	upload.Status = models.FAILED
	if err := v.VideosDAO.UpdateVideoTx(ctx, tx, video); err != nil {
		log.Errorf("Unable to update video with status  %v: %v", video.Status, err)
		return err
	}
	if err := v.UploadsDAO.UpdateUploadTx(ctx, tx, upload); err != nil {
		log.Errorf("Unable to update upload with status  %v: %v", upload.Status, err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Error("Cannot commit database transaction")
		return err
	}

	return nil
}

func writeHTTPResponse(video *models.Video, w http.ResponseWriter) {
	// Include videoCreated and status link into response (HATEOAS)
	links := map[string]jsonDTO.LinkJson{
		"status": jsonDTO.LinkToLinkJson(&models.Link{Href: "api/v1/videos/" + video.ID + "/status", Method: "GET"}),
		"stream": jsonDTO.LinkToLinkJson(&models.Link{Href: "api/v1/videos/" + video.ID + "/streams/master.m3u8", Method: "GET"}),
	}

	response := Response{
		Video: jsonDTO.VideoToVideoJson(video),
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

func (v VideoUploadHandler) publishStatus(video *models.Video) {
	msg, err := proto.Marshal(protobuf.VideoToVideoProtobuf(video))
	if err != nil {
		log.Error("Failed to Marshal status", err)
		return
	}

	if err := v.AmqpExchangerStatus.Publish(video.Title, msg); err != nil {
		log.Error("Unable to publish status update", err)
	}
}
