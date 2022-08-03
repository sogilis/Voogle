package controllers_test

import (
	"net/http/httptest"
	"testing"

	"github.com/Sogilis/Voogle/src/pkg/clients"

	"github.com/stretchr/testify/require"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
)

func TestHealthComponent(t *testing.T) { //nolint:cyclop
	cases := []struct {
		name             string
		expectedHTTPCode int
	}{
		{
			name:             "Success",
			expectedHTTPCode: 200,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			//Create request
			givenRequest := "/health"

			dummyServiceDiscovery := clients.NewDummyServiceDiscovery(nil, nil, nil, nil, nil)
			r := router.NewRouter(
				config.Config{},
				&router.Clients{ServiceDiscovery: dummyServiceDiscovery},
				&router.DAOs{},
			)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", givenRequest, nil)

			r.ServeHTTP(w, req)
			require.Equal(t, tt.expectedHTTPCode, w.Code)
		})
	}
}
