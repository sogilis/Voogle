package transformer

import (
	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/ffmpeg"
)

type GrayServer struct {
	TransformerServer
}

func newGrayServer(s3Client clients.IS3Client, discoveryClient clients.ServiceDiscovery) ITransformerServer {
	return &GrayServer{
		TransformerServer: TransformerServer{
			DiscoveryClient:         discoveryClient,
			S3Client:                s3Client,
			CreateTransformationCmd: ffmpeg.CreateGrayCommand,
		},
	}
}
