package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
)

type VideoInfo struct {
	Id    string `json:"id" example:"1"`
	Title string `json:"title" example:"my title"`
}
type AllVideos struct {
	Status string      `json:"status" example:"Success"`
	Data   []VideoInfo `json:"data"`
}

type VideosListHandler struct {
	MariadbClient *sql.DB
}

// VideosListHandler godoc
// @Summary Get list of all videos
// @Description Get list of all videos
// @Tags list
// @Accept  json
// @Produce  json
// @Success 200 {array} AllVideos
// @Failure 500 {object} object
// @Router /api/v1/videos/list [get]
func (v VideosListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("GET VideosListHandler")

	videos, err := dao.GetVideos(context.Background(), v.MariadbClient)
	if err != nil {
		log.Error("Unable to list objects from database: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	allVideos := AllVideos{}
	for _, video := range videos {
		videoInfo := VideoInfo{
			video.ID,
			video.Title,
		}
		allVideos.Data = append(allVideos.Data, videoInfo)
	}
	allVideos.Status = "Success"

	payload, err := json.Marshal(allVideos)

	if err != nil {
		log.Error("Unable to parse data struct in json ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(payload)
}
