package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Sogilis/Voogle/src/cmd/gray-server-transformer/config"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/ffmpeg"
	"github.com/Sogilis/Voogle/src/pkg/transformer/v1"
)

var _ transformer.TransformerServiceServer = &grayServer{}

type grayServer struct {
	transformer.UnimplementedTransformerServiceServer
	s3Client     clients.IS3Client
	ConsulClient clients.IConsulClient
}

func newGrayServer(s3 clients.IS3Client, consul clients.IConsulClient) *grayServer {
	return &grayServer{s3Client: s3, ConsulClient: consul}
}

func (r *grayServer) TransformVideo(ctx context.Context, args *transformer.TransformVideoRequest) (*transformer.TransformVideoResponse, error) {
	log.Debug("Beginning Transformation")

	var videoPart io.Reader
	if len(args.TransformerList) > 0 {
		log.Debug("Sending to next Service")

		// Select client for first tranformation and update list
		clientName := args.TransformerList[len(args.TransformerList)-1]
		args.TransformerList = args.TransformerList[:len(args.TransformerList)-1]

		// Retrieve service address
		tfService, err := r.ConsulClient.GetService(clientName)
		if err != nil {
			log.Errorf("Transformation service %v is unreachable %v ", clientName, err)
			return nil, err
		}

		// Create RPC client
		opts := grpc.WithTransportCredentials(insecure.NewCredentials())
		conn, err := grpc.Dial(tfService.Address+":"+tfService.Port, opts)
		if err != nil {
			log.Errorf("Cannot open TCP connection with grpc %v transformer server : %v", clientName, err)
			return nil, err
		}
		clientRPC := transformer.NewTransformerServiceClient(conn)

		// Transform video
		res, err := clientRPC.TransformVideo(ctx, args)
		if err != nil {
			log.Error("Failed to transform video", err)
			return nil, err
		}
		videoPart = bytes.NewReader(res.Data)
	} else {
		log.Debug("Retrieving on S3")

		// Retrieve the video part from bucket S3
		var err error
		videoPart, err = r.s3Client.GetObject(context.Background(), args.GetVideopath())
		if err != nil {
			log.Error("Failed to open video videoPath", err)
			return nil, err
		}
	}

	// Transform the video part
	transformedVideo, err := ffmpeg.TransformGrayscale(ctx, videoPart)
	if err != nil {
		log.Error("Cannot transformGrayscale")
		return nil, err
	}

	// Send transformed video part
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

	// ConsulClient to retrieve transformer address
	consulClient, err := clients.NewConsulClient(cfg.ConsulHost, cfg.ConsulUser, cfg.ConsulPwd)
	if err != nil {
		log.Fatal("Fail to create S3Client : ", err)
	}

	Addr := fmt.Sprintf("0.0.0.0:%v", cfg.Port)
	lis, err := net.Listen("tcp", Addr)
	if err != nil {
		log.Fatal("failed to listen : ", err)
	}

	grpcServer := grpc.NewServer()
	transformer.RegisterTransformerServiceServer(grpcServer, newGrayServer(s3Client, consulClient))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Cannot create gRPC server for Gray transformer service : ", err)
	}
}
