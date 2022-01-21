package controllers

import (
	"io"
	"net/http"

	cfg "github.com/Sogilis/Voogle/services/api/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type VideoGetMasterHandler struct {
	Cfg cfg.Config
}

func (v VideoGetMasterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debug("GET VideoGetMasterHandler - parameters ", vars)

	id, exist := vars["id"]
	if !exist {
		log.Error("Missing video id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Connection AWS S3 bucket
	creds := credentials.NewStaticCredentialsProvider(v.Cfg.S3AuthKey, v.Cfg.S3AuthPwd, "")

	cfg, err := config.LoadDefaultConfig(r.Context(), config.WithCredentialsProvider(creds), config.WithRegion(v.Cfg.S3Region))
	if err != nil {
		log.Error("Failed to create session S3 bucket", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	awsS3Client := s3.NewFromConfig(cfg)

	input := &s3.GetObjectInput{
		Bucket: aws.String(v.Cfg.S3Bucket),
		Key:    aws.String(id + "/master.m3u8"),
	}

	response, err := awsS3Client.GetObject(r.Context(), input)
	if err != nil {
		log.Error("Failed to open video "+id+"/master.m3u8 on S3 bucket "+v.Cfg.S3Bucket, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if _, err = io.Copy(w, response.Body); err != nil {
		log.Error("Unable to stream video master", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type VideoGetSubPartHandler struct {
	Cfg cfg.Config
}

func (v VideoGetSubPartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debug("GET VideoGetSubPartHandler - Parameters: ", vars)

	id, exist := vars["id"]
	if !exist {
		log.Error("Missing video id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	quality, exist := vars["quality"]
	if !exist {
		log.Error("Missing video quality")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filename, exist := vars["filename"]
	if !exist {
		log.Error("Missing video filename")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Connection AWS S3 bucket
	creds := credentials.NewStaticCredentialsProvider(v.Cfg.S3AuthKey, v.Cfg.S3AuthPwd, "")

	cfg, err := config.LoadDefaultConfig(r.Context(), config.WithCredentialsProvider(creds), config.WithRegion(v.Cfg.S3Region))
	if err != nil {
		log.Error("Failed to create session S3 bucket", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	awsS3Client := s3.NewFromConfig(cfg)

	input := &s3.GetObjectInput{
		Bucket: aws.String(v.Cfg.S3Bucket),
		Key:    aws.String(id + "/" + quality + "/" + filename),
	}

	response, err := awsS3Client.GetObject(r.Context(), input)
	if err != nil {
		log.Error("Failed to open video "+id+"/"+quality+"/"+filename+" on S3 bucket "+v.Cfg.S3Bucket, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if _, err := io.Copy(w, response.Body); err != nil {
		log.Error("Unable to stream subpart", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
