package controllers

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/cmd/api/metrics"
	"github.com/Sogilis/Voogle/src/pkg/clients"
	helpers "github.com/Sogilis/Voogle/src/pkg/transformer/helpers"
	"github.com/Sogilis/Voogle/src/pkg/transformer/v1"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"
)

type VideoGetMasterHandler struct {
	S3Client clients.IS3Client
	UUIDGen  uuidgenerator.IUUIDGenerator
}

// VideoGetMasterHandler godoc
// @Summary Get video master
// @Description Get video master
// @Tags video
// @Produce plain
// @Param id path string true "Video ID"
// @Success 200 {string} string "HLS video master"
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
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
	S3Client         clients.IS3Client
	UUIDGen          uuidgenerator.IUUIDGenerator
	ServiceDiscovery clients.ServiceDiscovery
}

// VideoGetSubPartHandler godoc
// @Summary Get sub part stream video
// @Description Get sub part stream video
// @Tags video
// @Produce plain
// @Param id path string true "Video ID"
// @Param quality path string true "Video quality"
// @Param filename path string true "Video sub part name"
// @Param filter query []string false "List of required filters"
// @Success 200 {string} string "Video sub part (.ts)"
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
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

	transformers := query["filter"]

	if strings.Contains(filename, "segment_index") || transformers == nil {
		object, err := v.S3Client.GetObject(r.Context(), id+"/"+quality+"/"+filename)
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
		// Add metrics (should be move into transformations service implem)
		for _, service := range transformers {
			if service == "gray" {
				metrics.CounterVideoTransformGray.Inc()
			} else if service == "flip" {
				metrics.CounterVideoTransformFlip.Inc()
			}
		}

		// Create transformation request
		request := transformer.TransformVideoRequest{
			Videopath:       id + "/" + quality + "/" + filename,
			TransformerList: transformers,
		}

		// Compute transformation(s) execution time for metrics
		start := time.Now()

		videoPart, err := helpers.GetVideoPart(r.Context(), &request, v.ServiceDiscovery, v.S3Client)
		if err != nil {
			log.Error("Unable to get video part", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		elapsed := time.Since(start)

		if len(transformers) == 1 {
			log.Debug("1 transform time : ", elapsed.Seconds())
			metrics.TransformationDuration.WithLabelValues("1").Observe(elapsed.Seconds())
		} else if len(transformers) == 2 {
			log.Debug("2 transform time : ", elapsed.Seconds())
			metrics.TransformationDuration.WithLabelValues("2").Observe(elapsed.Seconds())
		}

		if _, err := io.Copy(w, videoPart); err != nil {
			log.Error("Unable to stream subpart", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
