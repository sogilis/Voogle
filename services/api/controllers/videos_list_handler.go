package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	log "github.com/sirupsen/logrus"

	cfg "github.com/Sogilis/Voogle/services/api/config"
)

type VideoInfo struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}
type AllVideos struct {
	Status string      `json:"status"`
	Data   []VideoInfo `json:"data"`
}

type VideosListHandler struct {
	Cfg cfg.Config
}

func (v VideosListHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug("GET VideosListHandler")

	// Connection AWS S3 bucket
	creds := credentials.NewStaticCredentialsProvider(v.Cfg.S3AuthKey, v.Cfg.S3AuthPwd, "")

	cfg, err := config.LoadDefaultConfig(r.Context(), config.WithCredentialsProvider(creds), config.WithRegion(v.Cfg.S3Region))
	if err != nil {
		log.Error("Failed to create session S3 bucket", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	awsS3Client := s3.NewFromConfig(cfg)

	delimiter := "/"

	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(v.Cfg.S3Bucket),
		MaxKeys:   10,
		Delimiter: &delimiter,
	}
	res, err := awsS3Client.ListObjectsV2(r.Context(), input)
	if err != nil {
		log.Error("Failed to list video S3 bucket", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	allVideos := AllVideos{}
	for _, obj := range res.CommonPrefixes {

		if obj.Prefix == nil {
			continue
		}

		id := strings.TrimSuffix(*obj.Prefix, "/")

		log.Debug("S3 ID video:", id)
		videoInfo := VideoInfo{
			id,
			id,
		}
		allVideos.Data = append(allVideos.Data, videoInfo)
	}
	allVideos.Status = "Success"

	payload, err := json.Marshal(allVideos)

	if err != nil {
		log.Error("Unable to parse data struc in json", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(payload)
}
