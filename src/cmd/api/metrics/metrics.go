package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CounterApiVideoUploadInit = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_video_upload_init",
		Help: "The total number of processed events api video upload initialization",
	})
)

var (
	CounterApiVideoUploadSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_video_upload_sucess",
		Help: "The total number of processed events api video upload finish successfully",
	})
)
