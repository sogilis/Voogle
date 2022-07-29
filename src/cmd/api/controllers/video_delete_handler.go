package controllers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/pkg/clients"

	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	"github.com/Sogilis/Voogle/src/cmd/api/models"
)

type VideoDeleteHandler struct {
	S3Client   clients.IS3Client
	VideosDAO  *dao.VideosDAO
	UploadsDAO *dao.UploadsDAO
	UUIDGen    clients.IUUIDGenerator
}

// VideoDeleteHandler godoc
// @Summary Delete video
// @Description Delete video
// @Tags video
// @Produce plain
// @Param id path string true "Video ID"
// @Success 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/v1/videos/{id}/delete [delete]
func (v VideoDeleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debug("DELETE VideoDeleteHandler - parameters ", vars)

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
		log.Error("Cannot found video : ", err)
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if video.Status != models.ARCHIVE {
		log.Error("Video should be archived to be deleted", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	statusCode, err := v.deleteVideoAndUpload(r.Context(), id)
	if err != nil {
		w.WriteHeader(statusCode)
		return
	}

	if err = v.S3Client.RemoveObject(r.Context(), id); err != nil {
		log.Error("Cannot remove video "+id+" from S3 : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (v VideoDeleteHandler) deleteVideoAndUpload(ctx context.Context, id string) (int, error) {
	tx, err := v.VideosDAO.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Error("Cannot open new database transaction : ", err)
		return http.StatusInternalServerError, err
	}

	if err := v.UploadsDAO.DeleteUploadTx(ctx, tx, id); err != nil {
		log.Error("Cannot delete video "+id+" uploads : ", err)
		if err := tx.Rollback(); err != nil {
			log.Error("Cannot rollback : ", err)
		}

		return http.StatusInternalServerError, err
	}

	if err := v.VideosDAO.DeleteVideoTx(ctx, tx, id); err != nil {
		log.Error("Cannot delete video "+id+" : ", err)
		if err := tx.Rollback(); err != nil {
			log.Error("Cannot rollback : ", err)
		}

		return http.StatusInternalServerError, err
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
