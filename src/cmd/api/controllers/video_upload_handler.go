package controllers

import (
	"context"
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
func (v VideoUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("POST VideoUploadHandler")

	// Fetch video
	file, fileHandler, err := r.FormFile("video")
	if err != nil {
		log.Error("Missing file ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check if the received file is a supported video type
	if !isSupportedVideoType(file) {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	// Fetch cover
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

	// Fetch title
	title := r.FormValue("title")
	if title == "" {
		log.Error("Missing title file ")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Generate video UUID
	videoID, err := v.UUIDGen.GenerateUuid()
	if err != nil {
		log.Error("Cannot generate new videoID : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sourceName := "source" + filepath.Ext(fileHandler.Filename)
	sourcePath := videoID + "/" + sourceName

	coverName := ""
	coverPath := ""
	if fileCover != nil {
		coverName = "cover" + filepath.Ext(fileHandlerCover.Filename)
		coverPath = videoID + "/" + coverName
	}

	log.Infof("Receive video upload request with title : '%v'", title)

	// Create new video
	videoCreated, err := v.VideosDAO.CreateVideo(r.Context(), videoID, title, int(models.UPLOADING), sourcePath, coverPath)
	if err != nil {
		httpStatus := 0
		videoCreated, httpStatus = v.checkCreateVideo(r.Context(), w, title)

		if httpStatus != 0 {
			return
		}
	}

	metrics.CounterVideoUploadRequest.Inc()

	// Create new upload
	if err := v.uploadVideo(videoCreated, file, r); err != nil {
		metrics.CounterVideoUploadFail.Inc()
		log.Error("Cannot upload video : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	metrics.CounterVideoUploadSuccess.Inc()

	// Upload cover
	if err = v.uploadCover(videoCreated, fileCover, r); err != nil {
		log.Error("Cannot upload cover : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	metrics.CounterVideoEncodeRequest.Inc()
	if err = v.sendVideoForEncoding(r.Context(), videoCreated); err != nil {
		metrics.CounterVideoEncodeFail.Inc()
		log.Error("Cannot send video for encoding : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Include videoCreated and HATEOAS upload link into response
	writeHTTPResponse(videoCreated, w)
	log.Infof("Video '%v' successfully uploaded", title)
}

func isSupportedVideoType(input io.ReaderAt) bool {
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

func isSupportedCoverType(input io.ReaderAt) bool {
	supportedExtention := []string{".jpeg", ".jpg", ".png"}

	// Use ReadAt instead of Read to avoid seek affect resulting
	// in readed bytes missing
	buff := make([]byte, 262) // 262 bytes : no need more for image format
	if nbByte, err := input.ReadAt(buff, 0); err != nil {
		log.Error("Cannot check file type : ", err)
		log.Error("Asked 262 bytes, readed : ", nbByte)
		return false
	}

	mime := mimetype.Detect(buff)
	log.Debug("Receive " + mime.Extension() + " file")

	if !strings.EqualFold(mime.String()[:5], "image") {
		log.Error("Not supported file type : " + mime.String() + " (" + mime.Extension() + ")")
		return false
	}

	for _, ext := range supportedExtention {
		if strings.EqualFold(mime.Extension(), ext) {
			return true
		}
	}
	log.Error("Not supported file type : " + mime.String() + " (" + mime.Extension() + ")")
	return false
}

func (v VideoUploadHandler) checkCreateVideo(ctx context.Context, w http.ResponseWriter, title string) (*models.Video, int) {
	// Check if the returned error comes from duplicate title
	videoCreated, err := v.VideosDAO.GetVideoFromTitle(ctx, title)
	if err != nil {
		log.Error("Cannot find video ", title, "  : ", err)
		log.Error("Cannot insert new video into database: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return videoCreated, http.StatusInternalServerError
	}

	log.Debug("This title already exist, check video status")
	if videoCreated.Status == models.FAIL_UPLOAD {
		// Retry to upload+encode
		log.Debugf("Last upload of video %v failed, simply retry", videoCreated.Title)

	} else if videoCreated.Status == models.FAIL_ENCODE {
		// Retry to encode
		log.Debug("Ask for video encoding")
		metrics.CounterVideoEncodeRequest.Inc()
		if err = v.sendVideoForEncoding(ctx, videoCreated); err != nil {
			metrics.CounterVideoEncodeFail.Inc()
			log.Error("Cannot send video for encoding : ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return videoCreated, http.StatusInternalServerError
		}

		// Here, everything went well: video encoding asked and video status updated
		return videoCreated, http.StatusOK

	} else {
		// Title already exist, video already upload and encode, return error
		log.Error("A video with this title already uploaded and encoded")
		http.Error(w, "This title already exists", http.StatusConflict)
		return videoCreated, http.StatusConflict
	}
	return videoCreated, 0
}

func (v VideoUploadHandler) uploadVideo(video *models.Video, file multipart.File, r *http.Request) error {
	uploadID, err := v.UUIDGen.GenerateUuid()
	if err != nil {
		log.Error("Cannot generate new uploadID : ", err)
		return err
	}

	uploadCreated, err := v.UploadsDAO.CreateUpload(r.Context(), uploadID, video.ID, int(models.STARTED))
	if err != nil {

		// Update video status : FAIL_UPLOAD
		log.Error("Cannot insert new upload into database: ", err)
		video.Status = models.FAIL_UPLOAD
		if err := v.VideosDAO.UpdateVideo(r.Context(), video); err != nil {
			log.Errorf("Unable to update video with status  %v: %v", video.Status, err)
		}

		v.publishStatus(video)

		return err
	}

	// Upload video on S3
	err = v.S3Client.PutObjectInput(r.Context(), file, video.SourcePath)
	if err != nil {
		// Update video status : FAIL_UPLOAD + uploads status : FAILED
		log.Error("Unable to put object input on S3 ", err)

		if err := v.videoAndUploadFailed(r.Context(), video, uploadCreated); err != nil {
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
	if err = v.VideosDAO.UpdateVideo(r.Context(), video); err != nil {
		// Update video status : FAIL_UPLOAD + Update uploads status : FAILED
		log.Errorf("Unable to update video with status  %v : %v", video.Status, err)
		if err = v.videoAndUploadFailed(r.Context(), video, uploadCreated); err != nil {
			log.Error("video and upload status failed : ", err)
			return err
		}

		err = v.S3Client.RemoveObject(r.Context(), video.SourcePath)
		if err != nil {
			log.Errorf("Unable to remove uploaded video  %v : %v", video.ID, err)
		}
		return err
	}

	v.publishStatus(video)

	// Update uploads status : DONE + Upload date
	uploadCreated.Status = models.DONE
	uploadCreated.UploadedAt = &uploadDate
	if err = v.UploadsDAO.UpdateUpload(r.Context(), uploadCreated); err != nil {
		// Update uploads status : FAILED
		log.Errorf("Unable to update upload with status  %v: %v", uploadCreated.Status, err)
		if err = v.videoAndUploadFailed(r.Context(), video, uploadCreated); err != nil {
			log.Error("video and upload status failed : ", err)
			return err
		}

		err = v.S3Client.RemoveObject(r.Context(), video.SourcePath)
		if err != nil {
			log.Errorf("Unable to remove uploaded video  %v : %v", video.ID, err)
		}
	}

	return nil
}

func (v VideoUploadHandler) uploadCover(video *models.Video, file multipart.File, r *http.Request) error {
	if file == nil {
		return nil
	}

	err := v.S3Client.PutObjectInput(r.Context(), file, video.CoverPath)
	if err != nil {
		// Update video status : FAIL_UPLOAD + uploads status : FAILED
		log.Error("Unable to put object input on S3 ", err)

		return err
	}
	log.Debug("Success upload video cover " + video.ID + " on S3")
	return err
}

func (v VideoUploadHandler) sendVideoForEncoding(ctx context.Context, video *models.Video) error {
	videoProto := protobufDTO.VideoToVideoProtobuf(video)
	videoData, err := proto.Marshal(videoProto)
	if err != nil {
		log.Error("Unable to marshal video : ", err)
		v.videoEncodeFailed(ctx, video)
		return err
	}

	if err = v.AmqpClient.Publish(events.VideoUploaded, videoData); err != nil {
		log.Error("Unable to publish on Amqp client : ", err)
		v.videoEncodeFailed(ctx, video)
		return err
	}

	// Update video status : ENCODING
	video.Status = models.ENCODING
	if err := v.VideosDAO.UpdateVideo(ctx, video); err != nil {
		log.Errorf("Unable to update video with status  %v: %v", video.Status, err)
		v.videoEncodeFailed(ctx, video)
		return err
	}

	v.publishStatus(video)

	return nil
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

	video.Status = models.FAIL_UPLOAD
	upload.Status = models.FAILED
	if err := v.VideosDAO.UpdateVideoTx(ctx, tx, video); err != nil {
		log.Errorf("Unable to update video with status  %v: %v", video.Status, err)
		if err := tx.Rollback(); err != nil {
			log.Error("Cannot rollback : ", err)
		}
		return err
	}
	if err := v.UploadsDAO.UpdateUploadTx(ctx, tx, upload); err != nil {
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
		"status": jsonDTO.LinkToLinkJson(&models.Link{Href: "api/v1/videos/" + video.ID + "/status", Method: "GET"}),
		"stream": jsonDTO.LinkToLinkJson(&models.Link{Href: "api/v1/videos/" + video.ID + "/streams/master.m3u8", Method: "GET"}),
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

func (v VideoUploadHandler) publishStatus(video *models.Video) {
	msg, err := proto.Marshal(protobuf.VideoToVideoProtobuf(video))
	if err != nil {
		log.Error("Failed to Marshal status", err)
	}
	if err := v.AmqpExchangerStatus.Publish(video.Title, msg); err != nil {
		log.Error("Unable to publish status update", err)
	}
}
