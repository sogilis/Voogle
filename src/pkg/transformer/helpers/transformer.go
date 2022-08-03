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

const MAX_CHUNK_SIZE int = 32000

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
	var err error
	if len(args.TransformerList) > 0 {
		var transformedVideoPart bytes.Buffer
		err = sendToNextTransformer(ctx, args, &transformedVideoPart, serviceDiscovery)
		if err != nil {
			log.Error("Cannot send to next transformer : ", err)
			return nil, err
		}
		return &transformedVideoPart, nil
	}

	videoPart, err := getVideoFromS3(ctx, args.GetVideopath(), s3Client)
	if err != nil {
		log.Error("Cannot retrieve video from S3 : ", err)
		return nil, err
	}

	return videoPart, nil
}

func SendVideoPartStream(transformedVideoPartReader *io.PipeReader, stream transformer.TransformerService_TransformVideoServer) error {
	buf := make([]byte, MAX_CHUNK_SIZE)
	for {
		nbRead, err := transformedVideoPartReader.Read(buf)
		if err != nil {
			if err == io.EOF {
				stream.Context().Done()
				return nil
			}

			log.Error("Cannot transformFlip : ", err)
			_ = stream.Context().Err()
			return err
		}

		if nbRead != 0 {
			flipVideoPart := transformer.TransformVideoResponse{
				Chunk: buf[:nbRead],
			}

			if err := stream.Send(&flipVideoPart); err != nil {
				log.Error("Cannot send transformed video : ", err)
				_ = stream.Context().Err()
				return err
			}
		}
	}
}

func sendToNextTransformer(ctx context.Context, args *transformer.TransformVideoRequest, transformedVideoPart io.Writer, serviceDiscovery clients.ServiceDiscovery) error {
	// Select client for tranformation and update list
	clientName := args.TransformerList[len(args.TransformerList)-1]
	args.TransformerList = args.TransformerList[:len(args.TransformerList)-1]

	clientRPC, err := createRPCClient(clientName, serviceDiscovery)
	if err != nil {
		log.Errorf("Cannot create RPC Client %v : %v", clientName, err)
		return err
	}

	// Ask for next video part transformation
	streamResponse, err := clientRPC.TransformVideo(ctx, args)
	if err != nil {
		log.Error("Failed to transform video : ", err)
		return err
	}

	return recvVideoPartStream(transformedVideoPart, streamResponse)
}

func recvVideoPartStream(transformedVideoPart io.Writer, stream transformer.TransformerService_TransformVideoClient) error {
	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			log.Error("Failed to receive stream : ", err)
			return err
		}

		if res != nil {
			_, err := transformedVideoPart.Write(res.Chunk)
			if err != nil {
				log.Error("Failed to write : ", err)
				return err
			}
		}
	}
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
