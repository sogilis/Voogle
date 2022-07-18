package clients

import (
	"fmt"
	"strconv"
	"strings"

	consul_api "github.com/hashicorp/consul/api"
)

type ServiceDiscovery interface {
	GetTransformationServices(name string) ([]string, error)
	RegisterService(name, address string, port int, tags []string) error
}

var _ ServiceDiscovery = &serviceDiscovery{}

type serviceDiscovery struct {
	agent                     *consul_api.Agent
	transformersAddressesList map[string][]string
}

func NewServiceDiscovery(host, user, password string) (ServiceDiscovery, error) {
	config := &consul_api.Config{
		Address:  host,
		HttpAuth: &consul_api.HttpBasicAuth{Username: user, Password: password},
	}

	// Create a Consul API client
	client, err := consul_api.NewClient(config)
	if err != nil {
		return nil, err
	}

	// Create a Consul agent client
	agent := client.Agent()
	service := serviceDiscovery{agent: agent, transformersAddressesList: map[string][]string{}}
	if err := service.updateList(); err != nil {
		return nil, err
	}

	return &service, nil
}

// Get all available instances of transformation services
func (s *serviceDiscovery) updateList() error {
	services, err := s.agent.ServicesWithFilter("transformer in Service")
	if err != nil {
		return err
	}
	s.parseTransformerList(services)
	return nil
}

// Get all available instances of a given transformation service
func (s *serviceDiscovery) GetTransformationServices(name string) ([]string, error) {
	if len(s.transformersAddressesList[name]) < 1 {
		if err := s.updateList(); err != nil {
			return nil, err
		}
		if len(s.transformersAddressesList[name]) < 1 {
			return nil, fmt.Errorf("No service with name %v found.", name)
		}
	}
	return s.transformersAddressesList[name], nil
}

func (s *serviceDiscovery) parseTransformerList(services map[string]*consul_api.AgentService) {
	s.transformersAddressesList = map[string][]string{}
	for _, service := range services {
		name := strings.Split(service.Service, "-")[0]
		address := service.Address + ":" + strconv.Itoa(service.Port)
		s.transformersAddressesList[name] = append(s.transformersAddressesList[name], address)
	}
}

// Register service, will be use only for local dev
func (s *serviceDiscovery) RegisterService(name, address string, port int, tags []string) error {
	return s.agent.ServiceRegister(&consul_api.AgentServiceRegistration{
		Name:    name,
		Address: address,
		Port:    port,
		Tags:    tags,
	})
}
