package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Sogilis/Voogle/src/cmd/api/models"
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	"github.com/Sogilis/Voogle/src/cmd/api/models"
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
	vars := mux.Vars(r)

	attributeStr, exist := vars["attribute"]
	if !exist {
		log.Error("Missing sorting attribute")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var attribute models.PaginationAttribute
	switch attributeStr {
	case "title":
		attribute = models.TITLE
	case "upload_date":
		attribute = models.UPLOADEDAT
	case "creation_date":
		attribute = models.CREATEDAT
	case "update_date":
		attribute = models.UPDATEDAT
	default:
		log.Error("Attribute doesn't exist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	orderStr, exist := vars["order"]
	if !exist {
		log.Error("Missing sorting order")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	order, err := strconv.ParseBool(orderStr)
	if err != nil {
		log.Error("Order doesn't look like a boolean")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pageStr, exist := vars["page"]
	if !exist {
		log.Error("Missing page number")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		log.Error("Page not a number")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	limitStr, exist := vars["limit"]
	if !exist {
		log.Error("Missing limit number")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		log.Error("Limit not a number")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Debug("GET VideosListHandler")

	paginate := models.Pagination{
		Page:      uint(page),
		Limit:     uint(limit),
		Ascending: order,
		Attribute: attribute,
	}

	videos, err := dao.GetVideos(r.Context(), v.MariadbClient, paginate)
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
		if video.Status == models.COMPLETE {
			allVideos.Data = append(allVideos.Data, videoInfo)
		}
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
