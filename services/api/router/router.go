package router

import (
	"log"

	"github.com/Sogilis/Voogle/services/api/config"
	"github.com/Sogilis/Voogle/services/api/controllers"
	"github.com/goji/httpauth"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func NewRouter() (*mux.Router, *config.Config) {

	config, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to parse Env var", err)
	}

	r := mux.NewRouter()
	r.Use(httpauth.SimpleBasicAuth(config.UserAuth, config.PwdAuth))
	r.PathPrefix("/api/v1/videos").Handler(controllers.VideosListHandler{}).Methods("GET")

	return r, config
}

func GetCors() (handlers.CORSOption, handlers.CORSOption, handlers.CORSOption, handlers.CORSOption) {
	corsObj := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "OPTIONS"})
	headers := handlers.AllowedHeaders([]string{"Authorization"})
	credentials := handlers.AllowCredentials()

	return corsObj, methods, headers, credentials
}
