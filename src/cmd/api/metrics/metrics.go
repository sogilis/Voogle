package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var TotalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "api_http_requests_total",
		Help: "The total number of requests.",
	},
	[]string{"path"},
)

var ResponseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "api_response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var HttpDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "api_http_response_time_second",
		Help: "Duration of HTTP requests.",
	}, []string{"path"},
)

var (
	CounterVideoUploadSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_video_upload_sucess",
		Help: "The total number of processed events api video upload finish successfully",
	})
)