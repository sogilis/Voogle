package controllers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	"github.com/Sogilis/Voogle/src/cmd/api/models"
	"github.com/Sogilis/Voogle/src/pkg/clients"
)

type VideoUnarchiveHandler struct {
	VideosDAO *dao.VideosDAO
	UUIDGen   clients.IUUIDGenerator
}

// VideoUnarchiveHandler godoc
// @Summary Unarchive video
// @Description Unarchive video
// @Tags video
// @Produce plain
// @Param id path string true "Video ID"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/v1/videos/{id}/unarchive [put]
func (v VideoUnarchiveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debug("PUT VideoUnarchiveHandler - parameters ", vars)

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

	statusCode, err := v.unarchiveVideo(r.Context(), video)
	if err != nil {
		w.WriteHeader(statusCode)
		return
	}
}

func (v VideoUnarchiveHandler) unarchiveVideo(ctx context.Context, video *models.Video) (int, error) {
	// Can only unarchive video if it's in ARCHIVE state
	if video.Status != models.ARCHIVE {
		err := errors.New("Video status must be '" + models.ARCHIVE.String() + "' before getting '" + models.COMPLETE.String() + "'")
		log.Error(err)
		return http.StatusBadRequest, err
	}
	video.Status = models.COMPLETE

	if err := v.VideosDAO.UpdateVideo(ctx, video); err != nil {
		log.Error("Cannot update video "+video.ID+" : ", err)
		return http.StatusInternalServerError, err
	}
	return 0, nil
}
