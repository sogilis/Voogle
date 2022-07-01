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

type VideoArchiveVideoHandler struct {
	VideosDAO *dao.VideosDAO
	UUIDGen   uuidgenerator.IUUIDGenerator
}

// VideoArchiveVideoHandler godoc
// @Summary Archive video
// @Description Archive video
// @Tags status
// @Accept plain
// @Produce plain
// @Param id path string true "Video ID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /api/v1/videos/{id}/archive [put]
func (v VideoArchiveVideoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debug("GET VideoArchiveVideoHandler - parameters ", vars)

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

func (v VideoArchiveVideoHandler) archiveVideo(ctx context.Context, video *models.Video) (int, error) {
	tx, err := v.VideosDAO.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Error("Cannot open new database transaction : ", err)
		return http.StatusInternalServerError, err
	}

	// Can only archive video if it's in COMPLETE state
	if video.Status != models.COMPLETE {
		err := errors.New("Video status must be '" + models.COMPLETE.String() + "' before getting '" + models.ARCHIVE.String() + "'")
		log.Error(err)
		return http.StatusBadRequest, err
	}
	video.Status = models.ARCHIVE

	if err := v.VideosDAO.UpdateVideoTx(ctx, tx, video); err != nil {
		log.Error("Cannot update video "+video.ID+" : ", err)

		if err := tx.Rollback(); err != nil {
			log.Error("Cannot rollback : ", err)
			return http.StatusInternalServerError, err
		}

		if errors.Is(err, sql.ErrNoRows) {
			return http.StatusNotFound, err
		} else {
			return http.StatusInternalServerError, err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Error("Cannot commit database transaction")
		if err := tx.Rollback(); err != nil {
			log.Error("Cannot rollback : ", err)
		}
		return http.StatusInternalServerError, err
	}
	return 0, nil
}
