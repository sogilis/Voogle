package controllers

import (
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type VideoGetMasterHandler struct {
}

func (v VideoGetMasterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, exist := mux.Vars(r)["id"]
	if !exist {
		log.Error("Missing video id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	f, err := os.Open("./videos/" + id + "/master.m3u8")
	if err != nil {
		log.Error("Unable to open video master", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if _, err = io.Copy(w, f); err != nil {
		log.Error("Unable to stream video master", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
