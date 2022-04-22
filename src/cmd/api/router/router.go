package router

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/goji/httpauth"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/controllers"
	_ "github.com/Sogilis/Voogle/src/cmd/api/docs"
	"github.com/Sogilis/Voogle/src/cmd/api/metrics"
)

type Clients struct {
	S3Client      clients.IS3Client
	AmqpClient    clients.IAmqpClient
	MariadbClient *sql.DB
}

type UUIDGenerator struct {
	UUIDGen uuidgenerator.IUUIDGenerator
}
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// @title Orders API
// @version 1.0
// @description This is a sample serice for managing orders
// @termsOfService http://swagger.io/terms/
// @contact.name API Support Voogle
// @contact.email admin@sogilis.com
// @license.name AGPL
// @license.url LICENSE.txt
// @host localhost:4444
// @BasePath /
func NewRouter(config config.Config, clients *Clients, uuidGen *UUIDGenerator) http.Handler {
	metrics.InitMetrics()
	r := mux.NewRouter()
	r.Use(prometheusMiddleware)
	r.Use(httpauth.SimpleBasicAuth(config.UserAuth, config.PwdAuth))

	r.PathPrefix("/api/v1/videos/{id}/streams/master.m3u8").Handler(controllers.VideoGetMasterHandler{S3Client: clients.S3Client, UUIDGen: uuidGen.UUIDGen}).Methods("GET")
	r.PathPrefix("/api/v1/videos/{id}/streams/{quality}/{filename}").Handler(controllers.VideoGetSubPartHandler{S3Client: clients.S3Client, UUIDGen: uuidGen.UUIDGen}).Methods("GET")
	r.PathPrefix("/api/v1/videos/list").Handler(controllers.VideosListHandler{MariadbClient: clients.MariadbClient}).Methods("GET")
	r.PathPrefix("/api/v1/videos/upload").Handler(controllers.VideoUploadHandler{S3Client: clients.S3Client, AmqpClient: clients.AmqpClient, MariadbClient: clients.MariadbClient, UUIDGen: uuidGen.UUIDGen}).Methods("POST")
	r.PathPrefix("/api/v1/videos/{id}/status").Handler(controllers.VideoStatusHandler{MariadbClient: clients.MariadbClient, UUIDGen: uuidGen.UUIDGen}).Methods("GET")

	r.PathPrefix("/metrics").Handler(promhttp.Handler()).Methods("GET", "POST")
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	return handlers.CORS(getCORS())(r)
}

func getCORS() (handlers.CORSOption, handlers.CORSOption, handlers.CORSOption, handlers.CORSOption) {
	corsObj := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})
	headers := handlers.AllowedHeaders([]string{"Authorization"})
	credentials := handlers.AllowCredentials()

	return corsObj, methods, headers, credentials
}

// Metrics
func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(metrics.HttpDuration.WithLabelValues(path))
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		metrics.ResponseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		metrics.TotalRequests.WithLabelValues(path).Inc()

		timer.ObserveDuration()
	})
}
