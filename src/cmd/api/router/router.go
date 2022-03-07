package router

import (
	"net/http"
	"strconv"

	"github.com/goji/httpauth"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Sogilis/Voogle/src/pkg/clients"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/controllers"
	"github.com/Sogilis/Voogle/src/cmd/api/metrics"
)

type Clients struct {
	S3Client   clients.IS3Client
	AmqpClient clients.IAmqpClient
}
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewRouter(config config.Config, clients *Clients) http.Handler {
	metrics.InitMetrics()
	r := mux.NewRouter()
	r.Use(prometheusMiddleware)
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
