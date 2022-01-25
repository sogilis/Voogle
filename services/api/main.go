package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/services/api/clients"
	. "github.com/Sogilis/Voogle/services/api/config"
	"github.com/Sogilis/Voogle/services/api/router"
)

func main() {
	log.Info("Starting Voogle API")

	config, err := NewConfig()
	if err != nil {
		log.Fatal("Failed to parse Env var ", err)
	}
	if config.IsDev {
		log.SetLevel(log.DebugLevel)
	}

	s3Client, err := clients.NewS3Client(config.S3Region, config.S3Bucket, config.S3AuthKey, config.S3AuthPwd)
	if err != nil {
		log.Fatal("Failed to create S3 client - ", err)
	}

	routerClients := &router.Clients{
		S3Client: s3Client,
	}

	log.Info("Starting server on port:", config.Port)
	srv := &http.Server{
		Handler: router.NewRouter(config, routerClients),
		Addr:    fmt.Sprintf("0.0.0.0:%v", config.Port),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Crashed with error: ", err)
	}

}
