package metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var TotalRequests = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "api_http_requests_total",
		Help: "The total number of requests.",
	},
	[]string{"path"},
)

var ResponseStatus = promauto.NewCounterVec(
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

var TransformationDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "api_transformation_duration",
		Help: "Duration of video transformation request",
	},
	[]string{"transformations"},
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

func StoreTranformationTime(start time.Time, transformers []string) {
	elapsed := time.Since(start)
	if len(transformers) == 1 {
		TransformationDuration.WithLabelValues(transformers[0]).Observe(elapsed.Seconds())
	} else if len(transformers) > 1 {
		TransformationDuration.WithLabelValues(fmt.Sprint(len(transformers))).Observe(elapsed.Seconds())
	}
}
