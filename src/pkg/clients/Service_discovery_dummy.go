package clients

import (
	"github.com/Sogilis/Voogle/src/cmd/api/models"
)

var _ ServiceDiscovery = &dummyServiceDiscovery{}

type dummyServiceDiscovery struct {
	transformersAddressesList map[string]*TransformersInstances
	getTransformationService  func(s string) (string, error)
	getExistingServices       func(map[string]*TransformersInstances) []models.TransformerService
	startServiceDiscovery     func(serviceInfos ServiceInfos) error
	stop                      func()
}

func NewDummyServiceDiscovery(
	transformersAddressesList map[string]*TransformersInstances,
	getTransformationService func(s string) (string, error),
	getExistingServices func(map[string]*TransformersInstances) []models.TransformerService,
	startServiceDiscovery func(serviceInfos ServiceInfos) error,
	stop func(),
) ServiceDiscovery {
	return dummyServiceDiscovery{
		transformersAddressesList,
		getTransformationService,
		getExistingServices,
		startServiceDiscovery,
		stop,
	}
}

func (d dummyServiceDiscovery) GetTransformationService(s string) (string, error) {
	return d.getTransformationService(s)
}

func (d dummyServiceDiscovery) GetExistingServices() []models.TransformerService {
	return d.getExistingServices(d.transformersAddressesList)
}

func (d dummyServiceDiscovery) StartServiceDiscovery(serviceInfos ServiceInfos) error {
	return d.startServiceDiscovery(serviceInfos)
}

func (d dummyServiceDiscovery) Stop() {
	d.stop()
}
