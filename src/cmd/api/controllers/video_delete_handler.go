package controllers

import (
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

	if err := dao.DeleteVideo(r.Context(), v.MariadbClient, id); err != nil {
		log.Error("Cannot delete video "+id+" : ", err)
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if err := dao.DeleteUpload(r.Context(), v.MariadbClient, id); err != nil {
		log.Error("Cannot delete video "+id+" uploads : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := v.S3Client.RemoveObject(r.Context())

}
