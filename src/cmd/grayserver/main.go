package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/Sogilis/Voogle/src/cmd/grayserver/config"
	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/ffmpeg"

	"github.com/Sogilis/Voogle/src/pkg/transformer/v1"
)

var _ transformer.TransformerServiceServer = &grayServer{}

type grayServer struct {
	transformer.UnimplementedTransformerServiceServer
	s3Client clients.IS3Client
}

func newGrayServer(s3 clients.IS3Client) *grayServer {
	return &grayServer{s3Client: s3}
}

func (r *grayServer) TransformVideo(ctx context.Context, args *transformer.TransformVideoRequest) (*transformer.TransformVideoResponse, error) {
	// Parse video path on s3
	pathParts := strings.Split(args.GetVideopath(), "/")
	inputFileName := pathParts[len(pathParts)-1]
	outputFileName := "out_" + pathParts[len(pathParts)-1]

	// Retrieve the video part from bucket S3
	object, err := r.s3Client.GetObject(context.Background(), args.GetVideopath())
	if err != nil {
		log.Error("Failed to open video videoPath", err)
		return nil, err
	}

	// Write the video part into local file
	buf, err := io.ReadAll(object)
	if err != nil {
		log.Error("Cannot read in file")
		return nil, err
	}
	err = os.WriteFile(inputFileName, buf, 0666)
	if err != nil {
		log.Error("Cannot WriteFile")
		return nil, err
	}
	defer os.Remove(inputFileName)

	// Transform the video part, write the result into local file
	transformedVideo, err := ffmpeg.TransformGrayscale(inputFileName, outputFileName)
	if err != nil {
		log.Error("Cannot transformGrayscale")
		return nil, err
	}
	defer os.Remove(outputFileName)

	// Send transformed video part
	grayVideoPart := transformer.TransformVideoResponse{
		Data: transformedVideo.Bytes(),
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

	Addr := fmt.Sprintf("0.0.0.0:%v", cfg.Port)
	lis, err := net.Listen("tcp", Addr)
	if err != nil {
		log.Fatal("failed to listen : ", err)
	}

	grpcServer := grpc.NewServer()
	transformer.RegisterTransformerServiceServer(grpcServer, newGrayServer(s3Client))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Cannot create gRPC server for Gray transformer service : ", err)
	}
}
