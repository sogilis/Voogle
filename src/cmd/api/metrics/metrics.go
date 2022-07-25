package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

func InitMetrics() {
	if err := prometheus.Register(TotalRequests); err != nil {
		log.Warning("Unable to register metrics.TotalRequests prometheus : ", err)
	}
	if err := prometheus.Register(ResponseStatus); err != nil {
		log.Warning("Unable to register metrics.ResponseStatus prometheus : ", err)
	}
}

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
	CounterVideoUploadRequest = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_video_upload_request",
		Help: "The total number of upload request",
	})
)

var (
	CounterVideoUploadSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_video_upload_success",
		Help: "The total number of processed events api video upload finish successfully",
	})
)

var (
	CounterVideoUploadFail = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_video_upload_fail",
		Help: "The total number of processed events api video upload finish error",
	})
)

var (
	CounterVideoEncodeRequest = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_video_encode_request",
		Help: "The total number of encode request",
	})
)

var (
	CounterVideoEncodeSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_video_encode_success",
		Help: "The total number of processed events encoder video encode finish successfully",
	})
)

var (
	CounterVideoEncodeFail = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_video_encode_fail",
		Help: "The total number of processed events encoder video encode finish with error",
	})
)

var (
	CounterVideoTransformGray = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_gray_transformation_request",
		Help: "The total number of gray transformation request",
	})
)

var (
	CounterVideoTransformFlip = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_flip_transformation_request",
		Help: "The total number of flip transformation request",
	})
)
