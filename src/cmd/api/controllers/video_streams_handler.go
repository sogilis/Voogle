package controllers

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Sogilis/Voogle/src/cmd/api/metrics"
	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/transformer/v1"
)

type VideoGetMasterHandler struct {
	S3Client clients.IS3Client
	UUIDGen  clients.IUUIDGenerator
}

// VideoGetMasterHandler godoc
// @Summary Get video master
// @Description Get video master
// @Tags video
// @Produce plain
// @Param id path string true "Video ID"
// @Success 200 {string} string "HLS video master"
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/v1/videos/{id}/streams/master.m3u8 [get]
func (v VideoGetMasterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Debug("GET VideoGetMasterHandler - parameters ", vars)

	id := vars["id"]
	if !v.UUIDGen.IsValidUUID(id) {
		log.Error("Invalid id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	object, err := v.S3Client.GetObject(r.Context(), id+"/master.m3u8")
	if err != nil {
		log.Error("Failed to open video "+id+"/master.m3u8 ", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if _, err = io.Copy(w, object); err != nil {
		log.Error("Unable to stream video master", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type VideoGetSubPartHandler struct {
	S3Client         clients.IS3Client
	UUIDGen          clients.IUUIDGenerator
	ServiceDiscovery clients.ServiceDiscovery
}

// VideoGetSubPartHandler godoc
// @Summary Get sub part stream video
// @Description Get sub part stream video
// @Tags video
// @Produce plain
// @Param id path string true "Video ID"
// @Param quality path string true "Video quality"
// @Param filename path string true "Video sub part name"
// @Param filter query []string false "List of required filters"
// @Success 200 {string} string "Video sub part (.ts)"
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /api/v1/videos/{id}/streams/{quality}/{filename} [get]
func (v VideoGetSubPartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	query := r.URL.Query()
	log.Debug("GET VideoGetSubPartHandler - Parameters: ", vars)

	id := vars["id"]
	if !v.UUIDGen.IsValidUUID(id) {
		log.Error("Invalid id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	quality := vars["quality"]
	filename := vars["filename"]
	transformers := query["filter"]
	s3VideoPath := id + "/" + quality + "/" + filename

	if strings.Contains(filename, "segment_index") || transformers == nil {
		object, err := v.S3Client.GetObject(r.Context(), id+"/"+quality+"/"+filename)
		if err != nil {
			log.Error("Failed to open video videoPath", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if _, err := io.Copy(w, object); err != nil {
			log.Error("Unable to stream subpart", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		// Add metrics (should be move into transformations service implem)
		for _, service := range transformers {
			if service == "gray" {
				metrics.CounterVideoTransformGray.Inc()
			} else if service == "flip" {
				metrics.CounterVideoTransformFlip.Inc()
			}
		}

		videoPart, err := v.getVideoPart(r.Context(), s3VideoPath, transformers)
		if err != nil {
			log.Error("Cannot get video part : ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(w, videoPart); err != nil {
			log.Error("Unable to stream subpart", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (v VideoGetSubPartHandler) getVideoPart(ctx context.Context, s3VideoPath string, transformers []string) (io.Reader, error) {
	if len(transformers) == 0 {
		// Retrieve the video part from aws S3
		var err error
		videoPart, err := v.S3Client.GetObject(ctx, s3VideoPath)
		if err != nil {
			log.Error("Failed to get video from S3 : ", err)
			return nil, err
		}
		return videoPart, nil

	} else {
		// Ask for video part transformation
		start := time.Now()

		// Connect to RPC Client
		clientRPC, err := v.connectClientRPC(transformers[len(transformers)-1])
		if err != nil {
			log.Error("Cannot connect to RPC client : ", err)
			return nil, err
		}

		// Ask RPC Client for video transformation
		request := transformer.TransformVideoRequest{
			Videopath:       s3VideoPath,
			TransformerList: transformers,
		}
		streamResponse, err := clientRPC.TransformVideo(ctx, &request)
		if err != nil {
			log.Error("Failed to transform video : ", err)
			return nil, err
		}

		var videoPart bytes.Buffer
		for {
			res, err := streamResponse.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Error("Failed to receive stream : ", err)
				return nil, err
			}

			if res != nil {
				_, err := videoPart.Write(res.Chunk)
				if err != nil {
					log.Error("Failed to write : ", err)
					return nil, err
				}
			}
		}

		log.Debug("transformation execution time : ", time.Since(start).Seconds())
		metrics.StoreTranformationTime(start, transformers)
		return &videoPart, nil
	}
}

func (v VideoGetSubPartHandler) connectClientRPC(clientName string) (transformer.TransformerServiceClient, error) {
	// Retrieve service address and port
	tfServices, err := v.ServiceDiscovery.GetTransformationService(clientName)
	if err != nil {
		log.Errorf("Cannot get address for service name %v : %v", clientName, err)
		return nil, err
	}

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(tfServices, opts)
	if err != nil {
		log.Errorf("Cannot open TCP connection with grpc %v transformer server : %v", clientName, err)
		return nil, err
	}

	return transformer.NewTransformerServiceClient(conn), nil
}
