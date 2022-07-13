package controllers

import (
	b64 "encoding/base64"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"
)

type VideoCoverHandler struct {
	S3Client  clients.IS3Client
	VideosDAO *dao.VideosDAO
	UUIDGen   uuidgenerator.IUUIDGenerator
}

// VideoCoverHandler godoc
// @Summary Get video cover image in base64
// @Description Get video cover image in base64
// @Tags video
// @Accept plain
// @Produce plain
// @Param id path string true "Video ID"
// @Success 200 {string} string "video cover image in base64"
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/v1/videos/{id}/cover [get]
func (v VideoCoverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debug("GET VideoCoverHandler - parameters ", vars)

	id, exist := vars["id"]
	if !exist {
		log.Error("Missing video id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !v.UUIDGen.IsValidUUID(id) {
		log.Error("Invalid id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Fetch video cover path from DB
	video, err := v.VideosDAO.GetVideo(r.Context(), id)
	if err != nil {
		log.Error("Failed to get video "+id+" info from DB: ", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Fetch cover image from S3
	object, err := v.S3Client.GetObject(r.Context(), video.CoverPath)
	if err != nil {
		log.Error("Failed to open video cover "+video.CoverPath+": ", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	rawObject, err := io.ReadAll(object)
	if err != nil {
		log.Error("Failed to convert to base64 video cover "+video.CoverPath+": ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b64Object := b64.StdEncoding.EncodeToString(rawObject)

	w.Header().Set("Content-Type", "text/plain")
	if _, err = w.Write([]byte(b64Object)); err != nil {
		log.Error("Unable to write video cover: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
