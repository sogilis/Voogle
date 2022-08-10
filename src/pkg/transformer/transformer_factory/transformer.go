package transformer

import (
	"context"
	"fmt"
	"io"
	"net"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/ffmpeg"
	"github.com/Sogilis/Voogle/src/pkg/transformer/v1"
)

const MAX_CHUNK_SIZE int = 32000

type ITransformerServer interface {
	StartRPCServer(ctx context.Context, srv transformer.TransformerServiceServer, port uint32) error
	TransformVideo(ctx context.Context, args *transformer.TransformVideoRequest, stream transformer.TransformerService_TransformVideoServer) error
	Stop()
}

type TransformerServer struct {
	CreateTransformationCmd func(ctx context.Context) *exec.Cmd
	DiscoveryClient         clients.ServiceDiscovery
	S3Client                clients.IS3Client
}

func (t TransformerServer) StartRPCServer(ctx context.Context, srv transformer.TransformerServiceServer, port uint32) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
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

func (t TransformerServer) TransformVideo(ctx context.Context, args *transformer.TransformVideoRequest, stream transformer.TransformerService_TransformVideoServer) error {
	// Transformer will start video transformation, remove itself from the list
	args.TransformerList = args.TransformerList[:len(args.TransformerList)-1]
	if len(args.TransformerList) == 0 {
		// Retrieve the video part from aws S3
		videoPart, err := t.S3Client.GetObject(ctx, args.GetVideopath())
		if err != nil {
			log.Error("Failed to open video on S3 : ", err)
			return err
		}

		transformedVideoPartReader, transformedVideoPartWriter := io.Pipe()
		go func() {
			defer transformedVideoPartWriter.Close()
			if err := ffmpeg.TransformHLSPart(t.CreateTransformationCmd(ctx), videoPart, transformedVideoPartWriter); err != nil {
				log.Error("Cannot run ffmpeg command : ", err)
			}
		}()

		return t.sendVideoPartStream(transformedVideoPartReader, stream)

	} else {
		// Ask next transformer for videoPart. We will receive it as stream
		videoPart, err := t.sendToNextTransformer(ctx, args)
		if err != nil {
			log.Error("Cannot send to next transformer : ", err)
			return err
		}

		// Create transformation command and init a pipe for stdin
		cmd := t.CreateTransformationCmd(ctx)
		stdinWriter, err := cmd.StdinPipe()
		if err != nil {
			log.Error("Cannot create pipe stdin : ", err)
			return err
		}

		// Receive next transformer response, write it into stdin pipe
		go func() {
			defer stdinWriter.Close()
			err := t.recvVideoPartStream(stdinWriter, videoPart)
			if err != nil {
				log.Error("Failed to write : ", err)
				return
			}
		}()

		// Run the transformation command while we are receiving the file
		transformedVideoPartReader, transformedVideoPartWriter := io.Pipe()
		go func() {
			defer transformedVideoPartWriter.Close()

			// Execute command
			cmd.Stdout = transformedVideoPartWriter
			err = cmd.Start()
			if err != nil {
				log.Error("Cannot start command")
				return
			}

			// Wait end of ffmpeg command
			err = cmd.Wait()
			if err != nil {
				log.Error("Cannot wait command")
				return
			}
		}()

		return t.sendVideoPartStream(transformedVideoPartReader, stream)
	}
}

func (t TransformerServer) Stop() {
	t.DiscoveryClient.Stop()
}

func (t TransformerServer) sendToNextTransformer(ctx context.Context, args *transformer.TransformVideoRequest) (transformer.TransformerService_TransformVideoClient, error) {
	// Select client for tranformation and update list
	clientName := args.TransformerList[len(args.TransformerList)-1]

	clientRPC, err := t.createRPCClient(clientName)
	if err != nil {
		log.Errorf("Cannot create RPC Client %v : %v", clientName, err)
		return nil, err
	}

	// Ask for next video part transformation
	streamResponse, err := clientRPC.TransformVideo(ctx, args)
	if err != nil {
		log.Error("Failed to transform video : ", err)
		return nil, err
	}

	return streamResponse, nil
}

func (t TransformerServer) sendVideoPartStream(transformedVideoPartReader *io.PipeReader, stream transformer.TransformerService_TransformVideoServer) error {
	defer stream.Context().Done()

	buf := make([]byte, MAX_CHUNK_SIZE)
	for {
		nbRead, err := transformedVideoPartReader.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			log.Error("Cannot read video part : ", err)
			return err
		}

		if nbRead != 0 {
			videoPart := transformer.TransformVideoResponse{
				Chunk: buf[:nbRead],
			}

			if err := stream.Send(&videoPart); err != nil {
				log.Error("Cannot send transformed video : ", err)
				return err
			}
		}
	}
}

func (t TransformerServer) recvVideoPartStream(transformedVideoPart io.Writer, stream transformer.TransformerService_TransformVideoClient) error {
	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
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
	return nil
}

func (t TransformerServer) createRPCClient(clientName string) (transformer.TransformerServiceClient, error) {
	// Retrieve service address and port
	tfServices, err := t.DiscoveryClient.GetTransformationService(clientName)
	if err != nil {
		log.Errorf("Cannot get address for service name %v : %v", clientName, err)
		return nil, err
	}

	conn, err := grpc.Dial(tfServices, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("Cannot open TCP connection with grpc %v transformer server : %v", clientName, err)
		return nil, err
	}

	return transformer.NewTransformerServiceClient(conn), nil
}
