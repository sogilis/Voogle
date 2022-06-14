package clients

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Sogilis/Voogle/src/pkg/transformer/v1"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
)

type ITransformerManager interface {
	AddServiceClient(string, string) error
	TransformWithClients(context.Context, string, []string) (*transformer.TransformVideoResponse, error)
}

var _ ITransformerManager = &transformerManager{}

type transformerManager struct {
	s3Client IS3Client
	cfg      config.Config
	Clients  map[string]transformer.TransformerServiceClient
	adresses map[string]string
}

func NewTransformerManager(s3Client IS3Client, cfg config.Config) (ITransformerManager, error) {
	tsm := transformerManager{s3Client, cfg, map[string]transformer.TransformerServiceClient{}, map[string]string{}}

	return &tsm, nil
}

func (t *transformerManager) AddServiceClient(name string, adr string) error {
	// We connect to the service to ensure his status
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(adr, opts)
	if err != nil {
		log.Errorf("Cannot open TCP connection with grpc %v transformer server : %v", name, err)
		return err
	}
	client := transformer.NewTransformerServiceClient(conn)
	log.Debug("Client is connected :", name)
	// We send the existing name and adress to the new service
	err = sendAdressTo(client, t.adresses)
	if err != nil {
		log.Error("Could not send adress to service : ", err)
		return err
	}
	// We add the new service adress to the list and send it to the existing services
	for _, c := range t.Clients {
		err = sendAdressTo(c, map[string]string{name: adr})
		if err != nil {
			log.Error("Could not send adress to service : ", err)
			return err
		}
	}
	t.Clients[name] = client
	t.adresses[name] = adr
	return nil
}

// transformPath should be a list of name contained in the TransformerManager.Clients map and be at least 1 element long
func (t *transformerManager) TransformWithClients(c context.Context, videoPath string, transformList []string) (*transformer.TransformVideoResponse, error) {

	// We select the next client for current tranformation
	nextClientName := transformList[len(transformList)-1]
	nextClient := t.Clients[nextClientName]

	// We update the transformations left
	transformList = transformList[:len(transformList)-1]

	r := transformer.TransformVideoRequest{
		Videopath:                    videoPath,
		Additionnaltransformservices: transformList,
	}

	res, err := nextClient.TransformVideo(c, &r)
	if err != nil {
		log.Error("Could not transform video : ", err)
		return nil, err
	}

	return res, nil
}

func sendAdressTo(client transformer.TransformerServiceClient, adresses map[string]string) error {
	jsonadresses, err := json.Marshal(adresses)
	if err != nil {
		log.Error("Could not Marshal the adresses :", err)
		return err
	}
	r := transformer.SetTransformServicesRequest{
		Serviceadresses: jsonadresses,
	}
	_, err = client.SetTransformServices(context.Background(), &r)
	if err != nil {
		log.Error("Could not set service :", err)
		return err
	}
	return nil
}
