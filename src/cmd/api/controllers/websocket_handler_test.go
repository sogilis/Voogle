package controllers_test

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	hijack "github.com/getlantern/httptest"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
	"github.com/stretchr/testify/require"
)

func TestWebsocket(t *testing.T) { //nolint:cyclop
	requiredUsername := "valid"
	requiredPassword := "valid"

	cases := []struct {
		name             string
		givenUsername    string
		givenPassword    string
		expectedResponse int
	}{
		{
			name:             "Authentication Succeed",
			givenUsername:    "valid",
			givenPassword:    "valid",
			expectedResponse: 200,
		},
		{
			name:             "Authentication Fail with invalid Username",
			givenUsername:    "invalid",
			givenPassword:    "valid",
			expectedResponse: 403,
		},
		{
			name:             "Authentication Fail with invalid Password",
			givenUsername:    "valid",
			givenPassword:    "invalid",
			expectedResponse: 403,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			encodedAuth := "Basic%20" + base64.StdEncoding.EncodeToString([]byte(tt.givenUsername+":"+tt.givenPassword))
			authCookie := http.Cookie{Name: "Authorization", Value: encodedAuth}

			givenRequest := "/ws"

			r := router.NewRouter(config.Config{
				UserAuth: requiredUsername,
				PwdAuth:  requiredPassword,
			}, &router.Clients{}, &router.UUIDGenerator{}, &router.DAOs{})

			w := hijack.NewRecorder(nil)

			req := httptest.NewRequest("GET", givenRequest, nil)
			req.Header.Set("Connection", "Upgrade")
			req.Header.Set("Upgrade", "websocket")
			req.Header.Set("Sec-WebSocket-Version", "13")
			req.Header.Set("Sec-WebSocket-Key", "42")
			req.Header.Set("Sec-WebSocket-Extensions", "permessage-deflate")
			req.AddCookie(&authCookie)

			r.ServeHTTP(w, req)
			require.Equal(t, tt.expectedResponse, w.Code())

		})
	}

}
