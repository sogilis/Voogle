package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	jsonDTO "github.com/Sogilis/Voogle/src/cmd/api/dto/json"
	"github.com/Sogilis/Voogle/src/cmd/api/models"
)

type VideoInfo struct {
	Id    string `json:"id" example:"1"`
	Title string `json:"title" example:"my title"`
}

type VideoListResponse struct {
	Videos   []VideoInfo                 `json:"videos"`
	Links    map[string]jsonDTO.LinkJson `json:"_links"`
	LastPage int                         `json:"_lastpage"`
}

type VideosListHandler struct {
	MariadbClient *sql.DB
}

func checkRequest(vars map[string]string) (map[string]interface{}, error) {

	values := map[string]interface{}{}
	var err error

	//Check variables exists and are propers
	attributeStr, exist := vars["attribute"]
	if !exist {
		return values, errors.New("No sorting attribute")
	}

	switch attributeStr {
	case "title":
		values["attribute"] = models.TITLE
	case "upload_date":
		values["attribute"] = models.UPLOADEDAT
	case "creation_date":
		values["attribute"] = models.CREATEDAT
	case "update_date":
		values["attribute"] = models.UPDATEDAT
	default:
		return values, errors.New("Wrong attribute")
	}

	orderStr, exist := vars["order"]
	if !exist {
		return values, errors.New("No sorting order")
	}
	values["order"], err = strconv.ParseBool(orderStr)
	if err != nil {
		return values, errors.New("Doesn't look like a boolean")
	}

	pageStr, exist := vars["page"]
	if !exist {
		return values, errors.New("No page number")
	}
	values["page"], err = strconv.Atoi(pageStr)
	if err != nil {
		return values, errors.New("Page is not a number")
	}

	limitStr, exist := vars["limit"]
	if !exist {
		return values, errors.New("No limit number")
	}
	values["limit"], err = strconv.Atoi(limitStr)
	if err != nil {
		return values, errors.New("Limit is not a number")
	}
	return values, err
}

// VideosListHandler godoc
// @Summary Get list of all videos
// @Description Get list of all videos
// @Tags list
// @Accept  json
// @Produce  json
// @Success 200 {array} AllVideos
// @Failure 500 {object} object
// @Router /api/v1/videos/list/{attribute}/{order}/{page}/{limit} [get]
func (v VideosListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	vars, err := checkRequest(mux.Vars(r))
	if err != nil {
		log.Error("Request cannot be treated: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	attribute := vars["attribute"]
	order := vars["order"].(bool)
	page := vars["page"].(int)
	limit := vars["limit"].(int)

	log.Debug("GET VideosListHandler")

	//Initialize the response
	response := VideoListResponse{}

	//Get videos to be returned
	videos, err := dao.GetVideos(r.Context(), v.MariadbClient, attribute, order, page, limit)
	if err != nil {
		log.Error("Unable to list objects from database: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Add videos to response
	for _, video := range videos {
		if video.Status == models.COMPLETE {
			response.Videos = append(response.Videos, VideoInfo{
				Id:    video.ID,
				Title: video.Title,
			})
		}
	}

	//Add total number of page to the response
	totalvideos, err := dao.GetTotalVideos(r.Context(), v.MariadbClient)
	if err != nil {
		log.Error("Unable to get number of videos: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.LastPage = int(totalvideos / limit)
	if (totalvideos%limit) != 0 || response.LastPage == 0 {
		response.LastPage++
	}

	//Initialize path template and populate links response
	response.Links = map[string]jsonDTO.LinkJson{}
	path := "/api/v1/videos/list/" + mux.Vars(r)["attribute"] + "/" + mux.Vars(r)["order"] + "/%v/" + mux.Vars(r)["limit"]

	firstpath := fmt.Sprintf(path, 1)
	response.Links["first"] = jsonDTO.LinkToLinkJson(models.CreateLink(firstpath, "GET"))

	if page != 1 {
		previouspath := fmt.Sprintf(path, page-1)
		response.Links["previous"] = jsonDTO.LinkToLinkJson(models.CreateLink(previouspath, "GET"))
	}

	if page != response.LastPage {
		nextpath := fmt.Sprintf(path, page+1)
		response.Links["next"] = jsonDTO.LinkToLinkJson(models.CreateLink(nextpath, "GET"))
		lastpath := fmt.Sprintf(path, response.LastPage)
		response.Links["last"] = jsonDTO.LinkToLinkJson(models.CreateLink(lastpath, "GET"))
	}

	//Create and send the payload
	payload, err := json.Marshal(response)

	if err != nil {
		log.Error("Unable to parse data struct in json ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(payload)
}
