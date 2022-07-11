package clients

import (
	"strconv"

	"github.com/hashicorp/consul/api"
)

type TransformerInfos struct {
	Name    string
	Address string
	Port    string
}

type ServiceDiscovery interface {
	GetTransformationServices() ([]*TransformerInfos, error)
	GetTransformationServicesWithName(name string) ([]*TransformerInfos, error)
	RegisterService(name, address string, port int, tags []string) error
}

var _ ServiceDiscovery = &serviceDiscovery{}

type serviceDiscovery struct {
	agent *api.Agent
}

func NewServiceDiscovery(host, user, password string) (ServiceDiscovery, error) {
	config := &api.Config{
		Address:  host,
		HttpAuth: &api.HttpBasicAuth{Username: user, Password: password},
	}

	// Create a Consul API client
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	// Create a Consul agent client
	agent := client.Agent()
	service := serviceDiscovery{agent: agent}

	return service, nil
}

// Get all available instances of transformation services
func (s serviceDiscovery) GetTransformationServices() ([]*TransformerInfos, error) {
	services, err := s.agent.ServicesWithFilter("transformer in Tags")
	if err != nil {
		return nil, err
	}
	return s.parseTransformerInfos(services), nil
}

// Get all available instances of a given transformation service
func (s serviceDiscovery) GetTransformationServicesWithName(name string) ([]*TransformerInfos, error) {
	services, err := s.agent.ServicesWithFilter("transformer in Tags and " + name + " in Service")
	if err != nil {
		return nil, err
	}
	return s.parseTransformerInfos(services), nil
}

func (s serviceDiscovery) parseTransformerInfos(services map[string]*api.AgentService) []*TransformerInfos {
	transformers := make([]*TransformerInfos, 0, len(services))
	for _, service := range services {
		transformers = append(transformers, &TransformerInfos{
			Name:    service.Service,
			Address: service.Address,
			Port:    strconv.Itoa(service.Port),
		})
	}
	return transformers
}

// Register service, will be use only for local dev
func (s serviceDiscovery) RegisterService(name, address string, port int, tags []string) error {
	return s.agent.ServiceRegister(&api.AgentServiceRegistration{
		Name:    name,
		Address: address,
		Port:    port,
		Tags:    tags,
	})
}
