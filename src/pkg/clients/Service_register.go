package clients

import (
	"github.com/hashicorp/consul/api"
)

type ServiceRegister interface {
	RegisterService(name, address string, port int, tags []string) error
}

var _ ServiceRegister = &serviceRegister{}

type serviceRegister struct {
	agent *api.Agent
}

func NewServiceRegister(host, user, password string) (ServiceRegister, error) {
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
	service := serviceRegister{agent: agent}

	return service, nil
}

// Get all available instances of transformation services
func (s serviceRegister) RegisterService(name, address string, port int, tags []string) error {
	return s.agent.ServiceRegister(&api.AgentServiceRegistration{
		Name:    name,
		Address: address,
		Port:    port,
		Tags:    tags,
	})
}
