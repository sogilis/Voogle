package controllers_test

import (
	"net/http/httptest"
	"testing"

	"github.com/Sogilis/Voogle/src/pkg/clients"

	"github.com/stretchr/testify/require"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/models"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
)

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

	getExistingServices := func(list map[string]*clients.TransformersInstances) []models.TransformerService {
		existingServices := []models.TransformerService{}
		for name := range list {
			existingServices = append(existingServices, *models.CreateTransformerService(name))
		}
		return existingServices
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			//Create request
			givenRequest := "/api/v1/videos/transformer/list"

			dummyServiceDiscovery := clients.NewDummyServiceDiscovery(tt.adressCache, nil, getExistingServices, nil, nil)
			r := router.NewRouter(config.Config{
				UserAuth: givenUsername,
				PwdAuth:  givenPassword,
			}, &router.Clients{ServiceDiscovery: dummyServiceDiscovery}, &router.DAOs{})

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
