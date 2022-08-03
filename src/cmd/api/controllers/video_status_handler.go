package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	jsonDTO "github.com/Sogilis/Voogle/src/cmd/api/dto/json"
	"github.com/Sogilis/Voogle/src/pkg/clients"
)

type VideoGetStatusHandler struct {
	VideosDAO *dao.VideosDAO
	UUIDGen   clients.IUUIDGenerator
}

// VideoGetStatusHandler godoc
// @Summary Get video status
// @Description Get video status
// @Tags video
// @Produce plain
// @Param id path string true "Video ID"
// @Success 200 {object} jsonDTO.VideoStatus "Describe video status"
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/v1/videos/{id}/status [get]
func (v VideoGetStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debug("GET VideoGetStatusHandler - parameters ", vars)

	id := vars["id"]
	if !v.UUIDGen.IsValidUUID(id) {
		log.Error("Invalid id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	video, err := v.VideosDAO.GetVideo(r.Context(), id)
	if err != nil {
		log.Error("Cannot found video : ", err)
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	videoStatus := jsonDTO.VideoToStatusJson(video)
	payload, err := json.Marshal(videoStatus)
	if err != nil {
		log.Error("Unable to parse data struct in json ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(payload)
}
