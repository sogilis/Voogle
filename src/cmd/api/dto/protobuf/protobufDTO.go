package protobuf

import (
	log "github.com/sirupsen/logrus"

	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"

	"github.com/Sogilis/Voogle/src/cmd/api/models"
)

var protoToModelStatus = []models.VideoStatus{
	contracts.Video_VIDEO_STATUS_UNSPECIFIED: models.UNSPECIFIED,
	contracts.Video_VIDEO_STATUS_UPLOADING:   models.UPLOADING,
	contracts.Video_VIDEO_STATUS_UPLOADED:    models.UPLOADED,
	contracts.Video_VIDEO_STATUS_ENCODING:    models.ENCODING,
	contracts.Video_VIDEO_STATUS_COMPLETE:    models.COMPLETE,
	contracts.Video_VIDEO_STATUS_UNKNOWN:     models.UNKNOWN,
	contracts.Video_VIDEO_STATUS_FAIL_UPLOAD: models.FAIL_UPLOAD,
	contracts.Video_VIDEO_STATUS_FAIL_ENCODE: models.FAIL_ENCODE,
}

var modelToProtoStatus = []contracts.Video_VideoStatus{
	models.UNSPECIFIED: contracts.Video_VIDEO_STATUS_UNSPECIFIED,
	models.UPLOADING:   contracts.Video_VIDEO_STATUS_UPLOADING,
	models.UPLOADED:    contracts.Video_VIDEO_STATUS_UPLOADED,
	models.ENCODING:    contracts.Video_VIDEO_STATUS_ENCODING,
	models.COMPLETE:    contracts.Video_VIDEO_STATUS_COMPLETE,
	models.UNKNOWN:     contracts.Video_VIDEO_STATUS_UNKNOWN,
	models.FAIL_UPLOAD: contracts.Video_VIDEO_STATUS_FAIL_UPLOAD,
	models.FAIL_ENCODE: contracts.Video_VIDEO_STATUS_FAIL_ENCODE,
}

func VideoProtobufToVideo(videoProto *contracts.Video) *models.Video {
	if videoProto == nil {
		log.Error("Cannot convert protobuf video to video, video nil")
		return nil
	}

	video := models.Video{
		ID:         videoProto.Id,
		Title:      "",
		Status:     protoToModelStatus[videoProto.Status],
		UploadedAt: nil,
		CreatedAt:  nil,
		UpdatedAt:  nil,
	}

	return &video
}

func VideoToVideoProtobuf(video *models.Video, sourceName string) *contracts.Video {
	if video == nil {
		log.Error("Cannot convert protobuf video to video, video nil")
		return nil
	}

	videoData := &contracts.Video{
		Id:     video.ID,
		Status: modelToProtoStatus[video.Status],
		Source: sourceName,
	}

	return videoData
}
