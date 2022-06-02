package controllers

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"
)

type VideoGetMasterHandler struct {
	S3Client clients.IS3Client
	UUIDGen  uuidgenerator.IUUIDGenerator
}

// VideoGetMasterHandler godoc
// @Summary Get stream video
// @Description Get stream video
// @Tags streams
// @Accept plain
// @Produce plain
// @Param id path string true "Video ID"
// @Success 200 {string} string "OK"
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /api/v1/videos/{id}/streams/master.m3u8 [get]
func (v VideoGetMasterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debug("GET VideoGetMasterHandler - parameters ", vars)

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

	object, err := v.S3Client.GetObject(r.Context(), id+"/master.m3u8")
	if err != nil {
		log.Error("Failed to open video "+id+"/master.m3u8 ", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if _, err = io.Copy(w, object); err != nil {
		log.Error("Unable to stream video master", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type VideoGetSubPartHandler struct {
	S3Client           clients.IS3Client
	UUIDGen            uuidgenerator.IUUIDGenerator
	TransformerManager clients.ITransformerManager
}

// VideoGetSubPartHandler godoc
// @Summary Get sub part stream video
// @Description Get sub part stream video
// @Tags streams
// @Accept plain
// @Produce plain
// @Param id path string true "Video ID"
// @Param quality path string true "Video quality"
// @Param filename path string true "Video sub part name"
// @Success 200 {string} string "OK"
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /api/v1/videos/{id}/streams/{quality}/{filename} [get]
func (v VideoGetSubPartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	query := r.URL.Query()
	log.Debug("GET VideoGetSubPartHandler - Parameters: ", vars)

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

	quality, exist := vars["quality"]
	if !exist {
		log.Error("Missing video quality")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filename, exist := vars["filename"]
	if !exist {
		log.Error("Missing video filename")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if (strings.Contains(filename, "segment_index")) || (query["filter"] == nil) {
		object, err := v.S3Client.GetObject(context.Background(), id+"/"+quality+"/"+filename)
		if err != nil {
			log.Error("Failed to open video videoPath", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if _, err := io.Copy(w, object); err != nil {
			log.Error("Unable to stream subpart", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {

		videoPath := id + "/" + quality + "/" + filename
		videoPart, err := v.TransformerManager.TransformWithClients(r.Context(), videoPath, query["filter"])
		if err != nil {
			log.Error("Cannot transform video : ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(videoPart.Data); err != nil {
			log.Error("Unable to stream subpart", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
