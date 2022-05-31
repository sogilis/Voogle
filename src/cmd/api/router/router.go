package router

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"strconv"

	"github.com/goji/httpauth"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/controllers"
	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	_ "github.com/Sogilis/Voogle/src/cmd/api/docs"
	"github.com/Sogilis/Voogle/src/cmd/api/metrics"
)

type Clients struct {
	S3Client           clients.IS3Client
	AmqpClient         clients.IAmqpClient
	TransformerManager clients.ITransformerManager
}

type UUIDGenerator struct {
	UUIDGen uuidgenerator.IUUIDGenerator
}

type DAOs struct {
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
func NewRouter(config config.Config, clients *Clients, uuidGen *UUIDGenerator, DAOs *DAOs) http.Handler {
	metrics.InitMetrics()
	r := mux.NewRouter()
	r.Use(prometheusMiddleware)

	r.PathPrefix("/ws").Handler(controllers.WSHandler{Config: config}).Methods("GET")

	r.PathPrefix("/metrics").Handler(promhttp.Handler()).Methods("GET", "POST")
	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	v1 := r.PathPrefix("/api/v1").Subrouter()
	v1.Use(httpauth.SimpleBasicAuth(config.UserAuth, config.PwdAuth))

	v1.PathPrefix("/videos/{id}/streams/master.m3u8").Handler(controllers.VideoGetMasterHandler{S3Client: clients.S3Client, UUIDGen: uuidGen.UUIDGen}).Methods("GET")
	v1.PathPrefix("/videos/{id}/streams/{quality}/{filename}").Handler(controllers.VideoGetSubPartHandler{S3Client: clients.S3Client, UUIDGen: uuidGen.UUIDGen, TransformerManager: clients.TransformerManager}).Methods("GET")
	v1.PathPrefix("/videos/list/{attribute}/{order}/{page}/{limit}").Handler(controllers.VideosListHandler{VideosDAO: &DAOs.VideosDAO}).Methods("GET")
	v1.PathPrefix("/videos/{id}/delete").Handler(controllers.VideoDeleteVideoHandler{S3Client: clients.S3Client, VideosDAO: &DAOs.VideosDAO, UploadsDAO: &DAOs.UploadsDAO, UUIDGen: uuidGen.UUIDGen}).Methods("DELETE")
	v1.PathPrefix("/videos/{id}/info").Handler(controllers.VideoGetInfoHandler{VideosDAO: &DAOs.VideosDAO, UUIDGen: uuidGen.UUIDGen}).Methods("GET")
	v1.PathPrefix("/videos/list").Handler(controllers.VideosListHandler{VideosDAO: &DAOs.VideosDAO}).Methods("GET")
	v1.PathPrefix("/videos/upload").Handler(controllers.VideoUploadHandler{S3Client: clients.S3Client, AmqpClient: clients.AmqpClient, VideosDAO: &DAOs.VideosDAO, UploadsDAO: &DAOs.UploadsDAO, UUIDGen: uuidGen.UUIDGen}).Methods("POST")
	v1.PathPrefix("/videos/{id}/status").Handler(controllers.VideoGetStatusHandler{VideosDAO: &DAOs.VideosDAO, UUIDGen: uuidGen.UUIDGen}).Methods("GET")

	return handlers.CORS(getCORS())(r)
}

func getCORS() (handlers.CORSOption, handlers.CORSOption, handlers.CORSOption, handlers.CORSOption) {
	corsObj := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE"})
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

		metrics.ResponseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		metrics.TotalRequests.WithLabelValues(path).Inc()

		timer.ObserveDuration()
	})
}
