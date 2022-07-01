package clients

import (
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
)

type TransformerInfos struct {
	Name    string
	Address string
	Port    string
}

type ServiceDiscovery interface {
	GetTransformationServices() ([]TransformerInfos, error)
	GetTransformationServicesWithName(name string) ([]TransformerInfos, error)
}

var _ ServiceDiscovery = &serviceDiscovery{}

type serviceDiscovery struct {
	agent              *api.Agent
	localDockerNetwork bool
}

func NewServiceDiscovery(host, user, password string) (ServiceDiscovery, error) {
	config := &api.Config{
		Address:  host,
		HttpAuth: &api.HttpBasicAuth{Username: user, Password: password},
	}
	// Create a Consul API client
	client, err := api.NewClient(config)
	if err != nil {
		log.Error("Cannot create consul client")
	}
	// Create a Consul agent client
	agent := client.Agent()
	service := serviceDiscovery{agent: agent, localDockerNetwork: strings.Contains(host, "consul")}

	return service, nil
}

// Get all available instances of transformation services
func (s serviceDiscovery) GetTransformationServices() ([]TransformerInfos, error) {
	services, err := s.agent.ServicesWithFilter("transformer in Tags")
	if err != nil {
		log.Error("Cannot retrieve service :", err)
		return nil, err
	}
	return s.parseTransformerInfos(services), nil
}

// Get all available instances of a given transformation service
func (s serviceDiscovery) GetTransformationServicesWithName(name string) ([]TransformerInfos, error) {
	services, err := s.agent.ServicesWithFilter("transformer in Tags and " + name + " in Service")
	if err != nil {
		log.Error("Cannot retrieve service :", err)
		return nil, err
	}
	return s.parseTransformerInfos(services), nil
}

func (s serviceDiscovery) parseTransformerInfos(services map[string]*api.AgentService) []TransformerInfos {
	transformers := make([]TransformerInfos, 0, len(services))
	for _, service := range services {
		// If local dev is using docker network, use service name as address
		address := service.Address
		if s.localDockerNetwork {
			address = service.Service
		}
		transformers = append(transformers, TransformerInfos{
			Name:    service.Service,
			Address: address,
			Port:    strconv.Itoa(service.Port),
		})
	}
	return transformers
}
