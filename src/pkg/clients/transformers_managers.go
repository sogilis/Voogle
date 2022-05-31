package clients

import (
	"context"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/pkg/transformer/v1"
	log "github.com/sirupsen/logrus"
)

type ITransformerManager interface {
	AddServiceClient(string, transformer.TransformerServiceClient)
	TransformWithClient(context.Context, string, string) (*transformer.TransformVideoResponse, error)
}

var _ ITransformerManager = &transformerManager{}

type transformerManager struct {
	s3Client IS3Client
	cfg      config.Config
	Clients  map[string]transformer.TransformerServiceClient
}

func NewTransformerManager(s3Client IS3Client, cfg config.Config) (ITransformerManager, error) {
	tsm := transformerManager{s3Client, cfg, map[string]transformer.TransformerServiceClient{}}

	return &tsm, nil
}

func (t *transformerManager) AddServiceClient(name string, client transformer.TransformerServiceClient) {
	t.Clients[name] = client
}

func (t *transformerManager) TransformWithClient(c context.Context, name string, videoPath string) (*transformer.TransformVideoResponse, error) {
	r := transformer.TransformVideoRequest{
		Videopath: videoPath,
	}

	res, err := t.Clients[name].TransformVideo(c, &r)
	if err != nil {
		log.Error("Could not transform video : ", err)
		return nil, err
	}

	return res, nil
}
