package transformer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/transformer/v1"
)

func StartRPCServer(srv transformer.TransformerServiceServer, port uint32) {
	Addr := fmt.Sprintf("0.0.0.0:%v", port)
	lis, err := net.Listen("tcp", Addr)
	if err != nil {
		log.Fatal("failed to listen : ", err)
	}

	grpcServer := grpc.NewServer()
	transformer.RegisterTransformerServiceServer(grpcServer, srv)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Cannot create gRPC server : ", err)
	}
}

func GetVideoPart(ctx context.Context, args *transformer.TransformVideoRequest, serviceDiscovery clients.ServiceDiscovery, s3Client clients.IS3Client) (io.Reader, error) {
	var videoPart io.Reader
	var err error
	if len(args.TransformerList) > 0 {
		videoPart, err = sendToNextTransformer(ctx, args, serviceDiscovery)
		if err != nil {
			log.Error("Cannot send to next transformer : ", err)
			return nil, err
		}

	} else {
		videoPart, err = getVideoFromS3(ctx, args.GetVideopath(), s3Client)
		if err != nil {
			log.Error("Cannot retrieve video from S3 : ", err)
			return nil, err
		}
	}
	return videoPart, nil
}

func sendToNextTransformer(ctx context.Context, args *transformer.TransformVideoRequest, serviceDiscovery clients.ServiceDiscovery) (io.Reader, error) {
	// Select client for tranformation and update list
	clientName := args.TransformerList[len(args.TransformerList)-1]
	args.TransformerList = args.TransformerList[:len(args.TransformerList)-1]

	clientRPC, err := CreateRPCClient(clientName, serviceDiscovery)
	if err != nil {
		log.Errorf("Cannot create RPC Client %v : %v", clientName, err)
		return nil, err
	}

	// Ask for next video part transformation
	res, err := clientRPC.TransformVideo(ctx, args)
	if err != nil {
		log.Error("Failed to transform video : ", err)
		return nil, err
	}

	return bytes.NewReader(res.Data), nil
}

func getVideoFromS3(ctx context.Context, videoPath string, s3Client clients.IS3Client) (io.Reader, error) {
	// Retrieve the video part from aws S3
	object, err := s3Client.GetObject(ctx, videoPath)
	if err != nil {
		log.Error("Failed to open video videoPath : ", err)
		return nil, err
	}
	return object, nil
}

func CreateRPCClient(clientName string, serviceDiscovery clients.ServiceDiscovery) (transformer.TransformerServiceClient, error) {
	// Retrieve service address and port
	tfServices, err := serviceDiscovery.GetTransformationServicesWithName(clientName)
	if err != nil {
		log.Errorf("Transformation service %v is unreachable : %v ", clientName, err)
		return nil, err
	}

	if tfServices[0] != nil {
		// Create RPC client
		opts := grpc.WithTransportCredentials(insecure.NewCredentials())
		conn, err := grpc.Dial(tfServices[0].Address+":"+tfServices[0].Port, opts)
		if err != nil {
			log.Errorf("Cannot open TCP connection with grpc %v transformer server : %v", clientName, err)
			return nil, err
		}
		return transformer.NewTransformerServiceClient(conn), nil
	} else {
		return nil, fmt.Errorf("Service %v not found", clientName)
	}
}
