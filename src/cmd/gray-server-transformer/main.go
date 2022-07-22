package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/ffmpeg"
	helpers "github.com/Sogilis/Voogle/src/pkg/transformer/helpers"
	"github.com/Sogilis/Voogle/src/pkg/transformer/v1"

	"github.com/Sogilis/Voogle/src/cmd/gray-server-transformer/config"
)

const GOROUTINE_FLUSH_TIMEOUT time.Duration = time.Millisecond * 100

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
		log.Error("Cannot transformGrayscale : ", err)
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

	// serviceDiscovery to retrieve transformer address
	discoveryClient, err := clients.NewServiceDiscovery(cfg.ConsulHost, cfg.ConsulUser, cfg.ConsulPwd)
	if err != nil {
		log.Fatal("Fail to create Service Discovery : ", err)
	}

	// Start service discovery
	go func() {
		serviceInfos := clients.ServiceInfos{
			Name:    "gray-server-transformer",
			Address: cfg.LocalAddr,
			Port:    int(cfg.Port),
			Tags:    []string{"transformer"},
		}
		if err := discoveryClient.StartServiceDiscovery(serviceInfos); err != nil {
			log.Fatal("Discovery Service crash : ", err)
		}
	}()

	// Launch grpc Server
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		grayServer := &grayServer{
			s3Client:        s3Client,
			discoveryClient: discoveryClient,
		}
		if err := helpers.StartRPCServer(ctx, grayServer, cfg.Port); err != nil {
			log.Fatal("Gray RPC server error : ", err)
		}
	}()

	// Wait for SIGINT.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	// Stop discoveryClient and wait for grpcServer end properly
	cancel()
	discoveryClient.Stop()
	time.Sleep(GOROUTINE_FLUSH_TIMEOUT)
}
