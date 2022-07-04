package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"

	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	jsonDTO "github.com/Sogilis/Voogle/src/cmd/api/dto/json"
)

type VideoGetInfoHandler struct {
	VideosDAO *dao.VideosDAO
	UUIDGen   uuidgenerator.IUUIDGenerator
}

// VideoGetInfoHandler godoc
// @Summary Get video informations
// @Description Get video informations
// @Tags video
// @Produce json
// @Param id path string true "Video ID"
// @Success 200 {object} jsonDTO.VideoInfo "Video Informations"
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/v1/videos/{id}/info [get]
func (v VideoGetInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debug("GET VideoGetInfoHandler - parameters ", vars)

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

	videoInfo := jsonDTO.VideoToInfoJson(video)
	payload, err := json.Marshal(videoInfo)
	if err != nil {
		log.Error("Unable to parse data struct in json ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(payload)

}
