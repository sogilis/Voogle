package router

import (
	"net/http"

	"github.com/goji/httpauth"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Sogilis/Voogle/src/pkg/clients"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/controllers"
)

type Clients struct {
	S3Client   clients.IS3Client
	AmqpClient clients.IAmqpClient
}

func NewRouter(config config.Config, clients *Clients) http.Handler {
	r := mux.NewRouter()
	r.Use(httpauth.SimpleBasicAuth(config.UserAuth, config.PwdAuth))
	r.PathPrefix("/api/v1/videos/{id}/streams/master.m3u8").Handler(controllers.VideoGetMasterHandler{S3Client: clients.S3Client}).Methods("GET")
	r.PathPrefix("/api/v1/videos/{id}/streams/{quality}/{filename}").Handler(controllers.VideoGetSubPartHandler{S3Client: clients.S3Client}).Methods("GET")
	r.PathPrefix("/api/v1/videos/list").Handler(controllers.VideosListHandler{S3Client: clients.S3Client}).Methods("GET")
	r.PathPrefix("/api/v1/videos/upload").Handler(controllers.VideoUploadHandler{S3Client: clients.S3Client, AmqpClient: clients.AmqpClient}).Methods("POST")
	r.PathPrefix("/metrics").Handler(promhttp.Handler()).Methods("GET", "POST")

	return handlers.CORS(getCORS())(r)
}

func getCORS() (handlers.CORSOption, handlers.CORSOption, handlers.CORSOption, handlers.CORSOption) {
	corsObj := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})
	headers := handlers.AllowedHeaders([]string{"Authorization"})
	credentials := handlers.AllowCredentials()

	return corsObj, methods, headers, credentials
}
