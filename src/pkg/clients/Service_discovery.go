package clients

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	consul_api "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	log "github.com/sirupsen/logrus"
)

type ServiceDiscovery interface {
	GetTransformationServices(name string) ([]string, error)
	RegisterService(name, address string, port int, tags []string) error
	Watch(ctx context.Context, w chan string) error
}

var _ ServiceDiscovery = &serviceDiscovery{}

type serviceDiscovery struct {
	client                    *consul_api.Client
	agent                     *consul_api.Agent
	transformersAddressesList map[string][]string
}

func NewServiceDiscovery(consulURL, user, password string) (ServiceDiscovery, error) {
	config := &consul_api.Config{
		Address:  consulURL,
		HttpAuth: &consul_api.HttpBasicAuth{Username: user, Password: password},
	}

	// Create a Consul API client
	client, err := consul_api.NewClient(config)
	if err != nil {
		return nil, err
	}

	// Create service discovery
	service := serviceDiscovery{
		client:                    client,
		agent:                     client.Agent(),
		transformersAddressesList: map[string][]string{},
	}

	return &service, nil
}

// Get all available instances of a given transformation service
func (s *serviceDiscovery) GetTransformationServices(name string) ([]string, error) {
	if len(s.transformersAddressesList[name]) < 1 {
		return nil, fmt.Errorf("No service with name %v found.", name)
	}
	return s.transformersAddressesList[name], nil
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

// Get all available instances of transformation services
func (s *serviceDiscovery) updateList() error {
	services, err := s.agent.ServicesWithFilter("transformer in Service")
	if err != nil {
		return err
	}
	s.parseTransformerList(services)
	return nil
}

func (s *serviceDiscovery) parseTransformerList(services map[string]*consul_api.AgentService) {
	s.transformersAddressesList = map[string][]string{}
	for _, service := range services {
		name := strings.Split(service.Service, "-")[0]
		address := service.Address + ":" + strconv.Itoa(service.Port)
		s.transformersAddressesList[name] = append(s.transformersAddressesList[name], address)
	}
}

func (s *serviceDiscovery) Watch(ctx context.Context, w chan string) error {
	// Create Watch that check for changes on services register/deregister
	plan, err := watch.Parse(map[string]interface{}{"type": "services"})
	if err != nil {
		return err
	}

	// Define the handler function that will be called for each change
	plan.Handler = func(idx uint64, result interface{}) {
		log.Info("Change detected : Service register/deregister")
		log.Debug("index = ", idx, "\n", "result=", result)
		_ = s.updateList()
	}

	// Check for context
	go func(w chan string) {
		<-ctx.Done()
		log.Info("Gracefully shutdown consul watcher\n")
		plan.Stop()
		w <- "Closed"
	}(w)

	// Launch the watch. Note that the handler function will be run one first time
	err = plan.RunWithClientAndHclog(s.client, nil)

	// Should never be reached
	return err
}
