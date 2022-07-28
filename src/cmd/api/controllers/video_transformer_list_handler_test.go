package controllers_test

import (
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/Sogilis/Voogle/src/pkg/clients"

	"github.com/stretchr/testify/require"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
)

var _ clients.ServiceDiscovery = &dummyServiceDiscovery{}

type dummyServiceDiscovery struct {
	transformersAddressesList map[string]*clients.TransformersInstances
	mutex                     sync.RWMutex
}

func (d *dummyServiceDiscovery) GetTransformationService(string) (string, error) {
	return "", nil
}

func (d *dummyServiceDiscovery) StartServiceDiscovery(serviceInfos clients.ServiceInfos) error {
	return nil
}

func (d *dummyServiceDiscovery) Stop() {
}

func TestTransformerList(t *testing.T) { //nolint:cyclop

	//Initialize and set default parameters
	givenUsername := "dev"
	givenPassword := "test"
	cases := []struct {
		name             string
		authIsGiven      bool
		adressCache      map[string]*clients.TransformersInstances
		expectedHTTPCode int
	}{
		{
			name:             "Success",
			authIsGiven:      true,
			adressCache:      map[string]*clients.TransformersInstances{"s1": nil, "s2": nil},
			expectedHTTPCode: 200,
		},
		{
			name:             "Fail with no auth",
			authIsGiven:      false,
			adressCache:      map[string]*clients.TransformersInstances{},
			expectedHTTPCode: 401,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			//Create request
			givenRequest := "/api/v1/videos/transformer"

			dummyServiceDiscovery := &dummyServiceDiscovery{transformersAddressesList: tt.adressCache, mutex: sync.RWMutex{}}
			r := router.NewRouter(config.Config{
				UserAuth: givenUsername,
				PwdAuth:  givenPassword,
			}, &router.Clients{ServiceDiscovery: dummyServiceDiscovery}, &router.UUIDGenerator{}, &router.DAOs{})

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", givenRequest, nil)
			if tt.authIsGiven {
				req.SetBasicAuth(givenUsername, givenPassword)
			}

			r.ServeHTTP(w, req)
			require.Equal(t, tt.expectedHTTPCode, w.Code)
		})
	}
}
