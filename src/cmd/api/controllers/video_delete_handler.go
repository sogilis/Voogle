package controllers

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"

	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
)

type VideoDeleteVideoHandler struct {
	MariadbClient *sql.DB
	S3Client      clients.IS3Client
	UUIDGen       uuidgenerator.IUUIDGenerator
}

// VideoDeleteVideoHandler godoc
// @Summary Delete video
// @Description Delete video
// @Tags status
// @Accept plain
// @Produce plain
// @Param id path string true "Video ID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object

// @Router /api/v1/videos/{id}/delete [delete]
func (v VideoDeleteVideoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debug("DELETE VideoDeleteVideoHandler - parameters ", vars)

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

	video, err := dao.GetVideo(r.Context(), v.MariadbClient, id)
	if err != nil {
		log.Error("Cannot found video : ", err)
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	statusCode, err := v.deleteVideoAndUpload(r.Context(), video.ID)
	if err != nil {
		w.WriteHeader(statusCode)
		return
	}

	if err = v.S3Client.RemoveObject(r.Context(), video.SourcePath); err != nil {
		log.Error("Cannot remove video "+id+" from S3 : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (v VideoDeleteVideoHandler) deleteVideoAndUpload(ctx context.Context, id string) (int, error) {
	tx, err := v.MariadbClient.BeginTx(ctx, nil)
	if err != nil {
		log.Error("Cannot open new database transaction : ", err)
		return http.StatusInternalServerError, err
	}

	if err := dao.DeleteVideoTx(ctx, tx, id); err != nil {
		log.Error("Cannot delete video "+id+" : ", err)

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

	if err := dao.DeleteUploadTx(ctx, tx, id); err != nil {
		log.Error("Cannot delete video "+id+" uploads : ", err)
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
