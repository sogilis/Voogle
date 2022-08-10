package transformer

import (
	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/ffmpeg"
)

type FlipServer struct {
	TransformerServer
}

func newFlipServer(s3Client clients.IS3Client, discoveryClient clients.ServiceDiscovery) ITransformerServer {
	return &FlipServer{
		TransformerServer: TransformerServer{
			DiscoveryClient:         discoveryClient,
			S3Client:                s3Client,
			CreateTransformationCmd: ffmpeg.CreateFlipCommand,
		},
	}
}
