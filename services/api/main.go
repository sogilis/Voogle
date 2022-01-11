package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting Voogle API")

	r := mux.NewRouter()

	srv := &http.Server{
		Handler: r,
		Addr:    fmt.Sprintf("0.0.0.0:%v", 4444),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Crashed with error: ", err)
	}
}
