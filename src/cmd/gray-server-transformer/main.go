package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	transformer_factory "github.com/Sogilis/Voogle/src/pkg/transformer/transformer_factory"
	"github.com/Sogilis/Voogle/src/pkg/transformer/v1"

	"github.com/Sogilis/Voogle/src/cmd/gray-server-transformer/config"
)

const GOROUTINE_FLUSH_TIMEOUT time.Duration = time.Millisecond * 100

var _ transformer.TransformerServiceServer = &grayServer{}

type grayServer struct {
	transformer.UnimplementedTransformerServiceServer
	transformer transformer_factory.ITransformerServer
}

func (r *grayServer) TransformVideo(args *transformer.TransformVideoRequest, stream transformer.TransformerService_TransformVideoServer) error {
	log.Debug("Beginning Transformation")
	ctx := context.Background()
	return r.transformer.TransformVideo(ctx, args, stream)
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
	discoveryClient, err := clients.NewServiceDiscovery(cfg.ConsulHost)
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

	transformer, err := transformer_factory.GetTransformer("Gray", s3Client, discoveryClient)
	if err != nil {
		log.Error("Cannot create transformer : ", err)
	}

	// Launch grpc Server
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		grayServer := &grayServer{transformer: transformer}
		if err := transformer.StartRPCServer(ctx, grayServer, cfg.Port); err != nil {
			log.Fatal("Gray RPC server error : ", err)
		}
	}()

	// Wait for SIGINT.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	// Stop discoveryClient and wait for grpcServer end properly
	cancel()
	transformer.Stop()
	time.Sleep(GOROUTINE_FLUSH_TIMEOUT)
}
