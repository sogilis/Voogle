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

func StartRPCServer(ctx context.Context, srv transformer.TransformerServiceServer, port uint32) error {
	Addr := fmt.Sprintf("0.0.0.0:%v", port)
	lis, err := net.Listen("tcp", Addr)
	if err != nil {
		log.Error("failed to listen : ", err)
		return err
	}

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	// Check for context
	go func() {
		<-ctx.Done()
		log.Info("Gracefully shutdown grpcServer\n")
		grpcServer.Stop()
	}()

	transformer.RegisterTransformerServiceServer(grpcServer, srv)
	if err := grpcServer.Serve(lis); err != nil {
		log.Error("Cannot create gRPC server : ", err)
		return err
	}
	return nil
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

	clientRPC, err := createRPCClient(clientName, serviceDiscovery)
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

func createRPCClient(clientName string, serviceDiscovery clients.ServiceDiscovery) (transformer.TransformerServiceClient, error) {
	// Retrieve service address and port
	tfServices, err := serviceDiscovery.GetTransformationService(clientName)
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
