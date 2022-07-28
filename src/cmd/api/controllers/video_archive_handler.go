package controllers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"

	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	"github.com/Sogilis/Voogle/src/cmd/api/models"
)

type VideoArchiveHandler struct {
	VideosDAO *dao.VideosDAO
	UUIDGen   uuidgenerator.IUUIDGenerator
}

// VideoArchiveHandler godoc
// @Summary Archive video
// @Description Archive video
// @Tags video
// @Produce plain
// @Param id path string true "Video ID"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/v1/videos/{id}/archive [put]
func (v VideoArchiveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debug("GET VideoArchiveHandler - parameters ", vars)

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

	statusCode, err := v.archiveVideo(r.Context(), video)
	if err != nil {
		w.WriteHeader(statusCode)
		return
	}
}

func (v VideoArchiveHandler) archiveVideo(ctx context.Context, video *models.Video) (int, error) {
	// Can only archive video if it's in COMPLETE state
	if video.Status != models.COMPLETE {
		err := errors.New("Video status must be '" + models.COMPLETE.String() + "' before getting '" + models.ARCHIVE.String() + "'")
		log.Error(err)
		return http.StatusBadRequest, err
	}
	video.Status = models.ARCHIVE

	if err := v.VideosDAO.UpdateVideo(ctx, video); err != nil {
		log.Error("Cannot update video "+video.ID+" : ", err)
		return http.StatusInternalServerError, err
	}
	return 0, nil
}
