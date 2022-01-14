package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/services/api/config"
	"github.com/Sogilis/Voogle/services/api/router"
)

func main() {
	log.Info("Starting Voogle API")

	config, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to parse Env var", err)
	}

	log.Info("Starting server on port:", config.Port)
	srv := &http.Server{
		Handler: router.NewRouter(*config),
		Addr:    fmt.Sprintf("0.0.0.0:%v", config.Port),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Crashed with error: ", err)
	}

}
