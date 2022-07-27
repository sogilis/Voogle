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

type TransformersInstances struct {
	index        int
	servicesURLs []string
}

type ServiceDiscovery interface {
	GetTransformationService(name string) (string, error)
	StartServiceDiscovery(serviceInfos ServiceInfos) error
	Stop()
}

var _ ServiceDiscovery = &serviceDiscovery{}

type serviceDiscovery struct {
	client                    *consul_api.Client
	agent                     *consul_api.Agent
	plan                      *watch.Plan
	transformersAddressesList map[string]*TransformersInstances
	mutex                     sync.RWMutex
}

func NewServiceDiscovery(consulURL string) (ServiceDiscovery, error) {
	// Create a Consul API client
	client, err := consul_api.NewClient(&consul_api.Config{Address: consulURL})
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
		transformersAddressesList: map[string]*TransformersInstances{},
		mutex:                     sync.RWMutex{},
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
func (s *serviceDiscovery) GetTransformationService(name string) (string, error) {
	// We need to ensure that the Watch function runs by another goroutine is not
	// currently modifying the list
	s.mutex.RLock()
	if len(s.transformersAddressesList[name].servicesURLs) < 1 {
		return "", fmt.Errorf("No service with name %v found.", name)
	}
	serviceInstance := loadBalancing(s.transformersAddressesList[name])
	s.mutex.RUnlock()
	return serviceInstance, nil
}

func loadBalancing(t *TransformersInstances) string {
	t.index = t.index + 1
	if t.index >= len(t.servicesURLs) {
		t.index = 0
	}
	return t.servicesURLs[t.index]
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

	tmpList := map[string]*TransformersInstances{}

	for _, service := range services {
		name := strings.Split(service.Service, "-")[0]

		if tmpList[name] == nil {
			if s.transformersAddressesList[name] == nil {
				tmpList[name] = &TransformersInstances{index: 0, servicesURLs: []string{}}
			} else {
				tmpList[name] = &TransformersInstances{index: s.transformersAddressesList[name].index, servicesURLs: []string{}}
			}
		}
		// address := service.Address + ":" + strconv.Itoa(service.Port)

		// TODO : REMOVE IT WHEN SQUARESCALE UPDATE PORT FROM DOCKERFILE TO NOMAD/CONSUL
		address := service.Address + ":" + fixMeLater(name)
		tmpList[name].servicesURLs = append(tmpList[name].servicesURLs, address)
	}
	s.transformersAddressesList = tmpList
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
