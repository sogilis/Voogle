package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
)

type VideoStatus struct {
	Status string `json:"status" example:"UPLOADED"`
}

type VideoStatusHandler struct {
	MariadbClient *sql.DB
}

// @Summary Get video status
// @Description Get video status
// @Tags status
// @Accept plain
// @Produce plain
// @Param id path string true "Video ID"
// @Success 200 {Json} Json status:"Video status"
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router api/v1/videos/{id}/status [get]
func (v VideoStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debug("GET VideoGetMasterHandler - parameters ", vars)

	id, exist := vars["id"]
	if !exist {
		log.Error("Missing video id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	video, err := dao.GetVideo(v.MariadbClient, id)
	if err != nil {
		log.Error("Cannot get video from database : ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	videoStatus := VideoStatus{
		Status: video.Status.String(),
	}

	payload, err := json.Marshal(videoStatus)
	if err != nil {
		log.Error("Unable to parse data struct in json ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(payload)

}
