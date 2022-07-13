package main

import (
	"context"
	"io"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/ffmpeg"
	helpers "github.com/Sogilis/Voogle/src/pkg/transformer/helpers"
	"github.com/Sogilis/Voogle/src/pkg/transformer/v1"

	"github.com/Sogilis/Voogle/src/cmd/flip-server-transformer/config"
)

var _ transformer.TransformerServiceServer = &flipServer{}

type flipServer struct {
	transformer.UnimplementedTransformerServiceServer
	s3Client        clients.IS3Client
	discoveryClient clients.ServiceDiscovery
}

func (r *flipServer) TransformVideo(ctx context.Context, args *transformer.TransformVideoRequest) (*transformer.TransformVideoResponse, error) {
	log.Debug("Beginning Transformation")

	videoPart, err := helpers.GetVideoPart(ctx, args, r.discoveryClient, r.s3Client)
	if err != nil {
		log.Error("Cannot get video part : ", err)
		return nil, err
	}

	res, err := transformVideo(ctx, videoPart)
	if err != nil {
		log.Error("Cannot get video part : ", err)
		return nil, err
	}

	return res, nil
}

func transformVideo(ctx context.Context, videoPart io.Reader) (*transformer.TransformVideoResponse, error) {
	// Transform the video part
	transformedVideo, err := ffmpeg.TransformFlip(ctx, videoPart)
	if err != nil {
		log.Error("Cannot transformFlip : ", err)
		return nil, err
	}

	flipVideoPart := transformer.TransformVideoResponse{
		Data: transformedVideo,
	}
	return &flipVideoPart, err
}

func main() {
	log.Info("Starting Voogle flip transformer")

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to parse Env var : ", err)
	}
	if cfg.DevMode {
		log.SetLevel(log.DebugLevel)
	}

	// S3 client to access the videos
	s3Client, err := clients.NewS3Client(cfg.S3Host, cfg.S3Region, cfg.S3Bucket, cfg.S3AuthKey, cfg.S3AuthPwd)
	if err != nil {
		log.Fatal("Fail to create S3Client : ", err)
	}

	// serviceDiscovery to retrieve transformer address
	discoveryClient, err := clients.NewServiceDiscovery(cfg.ConsulHost, cfg.ConsulUser, cfg.ConsulPwd)
	if err != nil {
		log.Fatal("Fail to create S3Client : ", err)
	}

	// Register service on consul (only for local env)
	registerService(cfg, discoveryClient)

	// Launc RPC server
	flipServer := &flipServer{
		s3Client:        s3Client,
		discoveryClient: discoveryClient,
	}
	helpers.StartRPCServer(flipServer, cfg.Port)
}

func registerService(cfg config.Config, discoveryClient clients.ServiceDiscovery) {
	if cfg.LocalAddr != "" {
		err := discoveryClient.RegisterService("flip-server-transformer", cfg.LocalAddr, int(cfg.Port), []string{"transformer"})
		if err != nil {
			log.Fatal("Fail to register servce : ", err)
		}
	}
}
