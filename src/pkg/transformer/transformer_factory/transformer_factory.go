package transformer

import (
	"fmt"

	"github.com/Sogilis/Voogle/src/pkg/clients"
)

func GetTransformer(transformerType string, s3Client clients.IS3Client, discoveryClient clients.ServiceDiscovery) (ITransformerServer, error) {
	if transformerType == "Flip" {
		return newFlipServer(s3Client, discoveryClient), nil
	}
	if transformerType == "Gray" {
		return newGrayServer(s3Client, discoveryClient), nil
	}
	return nil, fmt.Errorf("Unknown transformer")
}
