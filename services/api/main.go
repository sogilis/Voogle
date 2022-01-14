package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/services/api/router"
)

func main() {
	log.Info("Starting Voogle API")

	r, config := router.NewRouter()
	corsObj, methods, headers, credentials := router.GetCors()

	log.Info("Starting server on port:", config.Port)
	srv := &http.Server{
		Handler: handlers.CORS(corsObj, headers, methods, credentials)(r),
		Addr:    fmt.Sprintf("0.0.0.0:%v", config.Port),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Crashed with error: ", err)
	}

}
