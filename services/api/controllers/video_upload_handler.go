package controllers

import (
	"net/http"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/services/api/clients"
)

type VideoUploadHandler struct {
	S3Client clients.IS3Client
}

func (v VideoUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("POST VideoUploadHandler")

	file, fileHandler, err := r.FormFile("video")
	if err != nil {
		log.Error("Missing file ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	if title == "" {
		log.Error("Missing title file ")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	title = strings.ReplaceAll(title, " ", "_")

	err = v.S3Client.PutObjectInput(r.Context(), file, title+"/source."+filepath.Ext(fileHandler.Filename))
	if err != nil {
		log.Error("Unable to put object input on S3 ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Debug("Success upload video " + title + " on S3")
}
