package main

import (
	"context"
	"io"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/ffmpeg"
	helpers "github.com/Sogilis/Voogle/src/pkg/transformer/helpers"
	"github.com/Sogilis/Voogle/src/pkg/transformer/v1"

	"github.com/Sogilis/Voogle/src/cmd/gray-server-transformer/config"
)

var _ transformer.TransformerServiceServer = &grayServer{}

type grayServer struct {
	transformer.UnimplementedTransformerServiceServer
	s3Client        clients.IS3Client
	discoveryClient clients.ServiceDiscovery
}

func (r *grayServer) TransformVideo(ctx context.Context, args *transformer.TransformVideoRequest) (*transformer.TransformVideoResponse, error) {
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
	transformedVideo, err := ffmpeg.TransformGrayscale(ctx, videoPart)
	if err != nil {
		log.Error("Cannot transformGrayscale")
		return nil, err
	}

	grayVideoPart := transformer.TransformVideoResponse{
		Data: transformedVideo,
	}
	return &grayVideoPart, err
}

func main() {
	log.Info("Starting Voogle gray transformer")

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

	// Register service on consul (only for local env)
	if cfg.LocalAddr != "" {
		resgisterClient, err := clients.NewServiceRegister(cfg.ConsulHost, cfg.ConsulUser, cfg.ConsulPwd)
		if err != nil {
			log.Fatal("Fail to create S3Client : ", err)
		}

		err = resgisterClient.RegisterService("gray", cfg.LocalAddr, int(cfg.Port), []string{"transformer"})
		if err != nil {
			log.Fatal("Fail to create S3Client : ", err)
		}
	}

	// serviceDiscovery to retrieve transformer address
	discoveryClient, err := clients.NewServiceDiscovery(cfg.ConsulHost, cfg.ConsulUser, cfg.ConsulPwd)
	if err != nil {
		log.Fatal("Fail to create S3Client : ", err)
	}

	// Launc RPC server
	grayServer := &grayServer{
		s3Client:        s3Client,
		discoveryClient: discoveryClient,
	}
	helpers.StartRPCServer(grayServer, cfg.Port)
}
