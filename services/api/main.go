package main

import (
	"fmt"
	"net/http"

	"github.com/goji/httpauth"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting Voogle API")

	config, err := NewConfig()
	if err != nil {
		log.Fatal("Failed to parse Env var", err)
	}

	r := mux.NewRouter()
	r.Use(httpauth.SimpleBasicAuth(config.UserAuth, config.PwdAuth))
	r.PathPrefix("/api/v1/videos").Handler(videosListHandler{}).Methods("GET")

	log.Info("Starting server on port:", config.Port)
	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("0.0.0.0:%v", config.Port),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Crashed with error: ", err)
	}

}
