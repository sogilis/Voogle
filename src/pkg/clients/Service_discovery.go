package clients

import (
	"fmt"
	"strings"
	"sync"

	consul_api "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	log "github.com/sirupsen/logrus"
)

type ServiceInfos struct {
	Name    string
	Address string
	Port    int
	Tags    []string
}

type ServiceDiscovery interface {
	GetTransformationServices(name string) ([]string, error)
	StartServiceDiscovery(serviceInfos ServiceInfos) error
	Stop()
}

var _ ServiceDiscovery = &serviceDiscovery{}

type serviceDiscovery struct {
	client                    *consul_api.Client
	agent                     *consul_api.Agent
	plan                      *watch.Plan
	transformersAddressesList map[string][]string
	mutex                     sync.RWMutex
}

func NewServiceDiscovery(consulURL string) (ServiceDiscovery, error) {
	log.Debug(consulURL)
	config := &consul_api.Config{
		Address: consulURL,
	}

	// Create a Consul API client
	client, err := consul_api.NewClient(config)
	if err != nil {
		return nil, err
	}

	// Create Watch that check for changes on services register/deregister
	plan, err := watch.Parse(map[string]interface{}{"type": "services"})
	if err != nil {
		return nil, err
	}

	// Create service discovery
	service := serviceDiscovery{
		client:                    client,
		agent:                     client.Agent(),
		plan:                      plan,
		transformersAddressesList: map[string][]string{},
		mutex:                     sync.RWMutex{},
	}

	// DEBUG
	if log.GetLevel() == log.DebugLevel {
		services, err := service.agent.ServicesWithFilter("transformer in Service")
		if err != nil {
			log.Debug("consul request : ", err)
		}
		for _, service := range services {
			log.Debug("name: ", service.Service)
			log.Debug("address: ", service.Address)
		}
		catalog := client.Catalog()
		res, _, err := catalog.Services(nil)
		if err != nil {
			log.Debug("catalog : ", err)
		}
		for _, service := range res {
			log.Debug("res : ", service)
		}
	}

	return &service, nil
}

func (s *serviceDiscovery) StartServiceDiscovery(serviceInfos ServiceInfos) error {
	// Register service on consul (only for local env)
	if err := s.registerService(serviceInfos); err != nil {
		return err
	}

	// Run watcher on services, it updates the transformation service cache if a new service
	// is register/deregister on consul
	if err := s.watch(); err != nil {
		return err
	}

	return nil
}

// Get all available instances of a given transformation service
func (s *serviceDiscovery) GetTransformationServices(name string) ([]string, error) {
	// We need to ensure that the Watch function runs by another goroutine is not
	// currently modifying the list
	s.mutex.RLock()
	serviceInstances := s.transformersAddressesList[name]
	if len(serviceInstances) < 1 {
		return nil, fmt.Errorf("No service with name %v found.", name)
	}
	s.mutex.RUnlock()
	return serviceInstances, nil
}

func (s *serviceDiscovery) registerService(serviceInfos ServiceInfos) error {
	if serviceInfos.Address != "" {
		return s.agent.ServiceRegister(&consul_api.AgentServiceRegistration{
			Name:    serviceInfos.Name,
			Address: serviceInfos.Address,
			Port:    serviceInfos.Port,
			Tags:    serviceInfos.Tags,
		})
	}
	return nil
}

func (s *serviceDiscovery) updateList() error {
	services, err := s.agent.ServicesWithFilter("transformer in Service")
	if err != nil {
		return err
	}
	// This part update the transformersAddressesList, since it is down
	// by the Watch function which aims to be run by a goroutine, we need
	// to ensure that no one is already reading on this list.
	s.mutex.Lock()
	s.transformersAddressesList = map[string][]string{}
	for _, service := range services {
		name := strings.Split(service.Service, "-")[0]
		// address := service.Address + ":" + strconv.Itoa(service.Port)

		////////
		// TODO : REMOVE IT WHEN SQUARESCALE UPDATE PORT FROM DOCKERFILE TO NOMAD/CONSUL
		address := service.Address + ":" + fixMeLater(name)
		////////

		s.transformersAddressesList[name] = append(s.transformersAddressesList[name], address)
	}
	s.mutex.Unlock()
	return nil
}

func (s *serviceDiscovery) watch() error {
	// Define the handler function that will be called for each change
	s.plan.Handler = func(idx uint64, result interface{}) {
		log.Info("Change detected : Service register/deregister")
		log.Debug("index = ", idx, "\n", "result=", result)
		_ = s.updateList()
	}

	// Launch the watch. Note that the handler function will be run one first time
	return s.plan.RunWithClientAndHclog(s.client, nil)
}

func (s *serviceDiscovery) Stop() {
	s.plan.Stop()
	log.Info("Gracefully shutdown service discovery")
}

// TODO : REMOVE IT WHEN SQUARESCALE UPDATE PORT FROM DOCKERFILE TO NOMAD/CONSUL
func fixMeLater(name string) string {
	if name == "gray" {
		return "50051"
	} else if name == "flip" {
		return "50052"
	}
	return ""
}
