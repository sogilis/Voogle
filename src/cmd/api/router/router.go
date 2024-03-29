package router

import (
	"bufio"
	"database/sql"
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/goji/httpauth"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/Sogilis/Voogle/src/pkg/clients"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/controllers"
	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	_ "github.com/Sogilis/Voogle/src/cmd/api/docs"
	"github.com/Sogilis/Voogle/src/cmd/api/metrics"
)

type Clients struct {
	S3Client              clients.IS3Client
	AmqpClient            clients.AmqpClient
	AmqpVideoStatusUpdate clients.AmqpClient
	ServiceDiscovery      clients.ServiceDiscovery
	UUIDGen               clients.IUUIDGenerator
}
type DAOs struct {
	Db         *sql.DB
	VideosDAO  dao.VideosDAO
	UploadsDAO dao.UploadsDAO
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
func NewRouter(config config.Config, clients *Clients, DAOs *DAOs) http.Handler {
	r := mux.NewRouter()
	r.Use(prometheusMiddleware)

	r.PathPrefix("/ws").Handler(controllers.WSHandler{Config: config, AmqpVideoStatusUpdate: clients.AmqpVideoStatusUpdate}).Methods("GET")

	r.PathPrefix("/metrics").Handler(promhttp.Handler()).Methods("GET", "POST")
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	r.PathPrefix("/health").Handler(controllers.HealthComponentHandler{}).Methods("GET")

	v1 := r.PathPrefix("/api/v1").Subrouter()
	v1.Use(httpauth.SimpleBasicAuth(config.UserAuth, config.PwdAuth))

	v1.PathPrefix("/videos/{id}/streams/master.m3u8").Handler(controllers.VideoGetMasterHandler{S3Client: clients.S3Client, UUIDGen: clients.UUIDGen}).Methods("GET")
	v1.PathPrefix("/videos/{id}/streams/{quality}/{filename}").Handler(controllers.VideoGetSubPartHandler{S3Client: clients.S3Client, UUIDGen: clients.UUIDGen, ServiceDiscovery: clients.ServiceDiscovery}).Methods("GET")
	v1.PathPrefix("/videos/transformer/list").Handler(controllers.VideoTransformerListHandler{ServiceDiscovery: clients.ServiceDiscovery}).Methods("GET")
	v1.PathPrefix("/videos/{id}/cover").Handler(controllers.VideoCoverHandler{S3Client: clients.S3Client, VideosDAO: &DAOs.VideosDAO, UUIDGen: clients.UUIDGen}).Methods("GET")
	v1.PathPrefix("/videos/list/{attribute}/{order}/{page}/{limit}/{status}").Handler(controllers.VideosListHandler{VideosDAO: &DAOs.VideosDAO}).Methods("GET")
	v1.PathPrefix("/videos/{id}/delete").Handler(controllers.VideoDeleteHandler{S3Client: clients.S3Client, VideosDAO: &DAOs.VideosDAO, UploadsDAO: &DAOs.UploadsDAO, UUIDGen: clients.UUIDGen}).Methods("DELETE")
	v1.PathPrefix("/videos/{id}/archive").Handler(controllers.VideoArchiveHandler{VideosDAO: &DAOs.VideosDAO, UUIDGen: clients.UUIDGen}).Methods("PUT")
	v1.PathPrefix("/videos/{id}/unarchive").Handler(controllers.VideoUnarchiveHandler{VideosDAO: &DAOs.VideosDAO, UUIDGen: clients.UUIDGen}).Methods("PUT")
	v1.PathPrefix("/videos/{id}/info").Handler(controllers.VideoGetInfoHandler{VideosDAO: &DAOs.VideosDAO, UUIDGen: clients.UUIDGen}).Methods("GET")
	v1.PathPrefix("/videos/upload").Handler(controllers.VideoUploadHandler{S3Client: clients.S3Client, AmqpClient: clients.AmqpClient, AmqpVideoStatusUpdate: clients.AmqpVideoStatusUpdate, VideosDAO: &DAOs.VideosDAO, UploadsDAO: &DAOs.UploadsDAO, UUIDGen: clients.UUIDGen}).Methods("POST")
	v1.PathPrefix("/videos/{id}/status").Handler(controllers.VideoGetStatusHandler{VideosDAO: &DAOs.VideosDAO, UUIDGen: clients.UUIDGen}).Methods("GET")

	return handlers.CORS(getCORS())(r)
}

func getCORS() (handlers.CORSOption, handlers.CORSOption, handlers.CORSOption, handlers.CORSOption) {
	corsObj := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "OPTIONS", "DELETE"})
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

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(metrics.HttpDuration.WithLabelValues(path))
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		if strings.Contains(path, "/api/v1") {
			metrics.ResponseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
			metrics.TotalRequests.WithLabelValues(path).Inc()
		}

		timer.ObserveDuration()
	})
}
