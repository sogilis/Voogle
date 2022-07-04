package controllers_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
)

func TestVideoCover(t *testing.T) { //nolint:cyclop
	givenUsername := "dev"
	givenUserPwd := "test"

	validVideoID := "1508e7d5-5bc6-4a50-9176-ab0371aa65fe"
	invalidVideoID := "invalidvideoid"
	UUIDValidFunc := func(u string) bool { _, err := uuid.Parse(u); return err == nil }
	getObjectS3 := func(v string) (io.Reader, error) { return strings.NewReader(""), nil }

	cases := []struct {
		name             string
		giveRequest      string
		giveWithAuth     bool
		expectedHTTPCode int
		isValidUUID      func(string) bool
		getObject        func(string) (io.Reader, error)
	}{
		{
			name:             "GET video cover",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/cover",
			giveWithAuth:     true,
			expectedHTTPCode: 200,
			isValidUUID:      UUIDValidFunc,
			getObject:        getObjectS3,
		},
		{
			name:             "GET fails with invalid video ID",
			giveRequest:      "/api/v1/videos/" + invalidVideoID + "/cover",
			giveWithAuth:     true,
			expectedHTTPCode: 400,
			isValidUUID:      UUIDValidFunc,
			getObject:        getObjectS3,
		},
		{
			name:             "GET fails with no auth",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/cover",
			giveWithAuth:     false,
			expectedHTTPCode: 401,
			isValidUUID:      UUIDValidFunc,
			getObject:        getObjectS3,
		},
		{
			name:             "GET fails with S3 error",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/cover",
			giveWithAuth:     true,
			expectedHTTPCode: 404,
			isValidUUID:      UUIDValidFunc,
			getObject:        func(v string) (io.Reader, error) { return nil, fmt.Errorf("S3 error") },
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			s3Client := clients.NewS3ClientDummy(nil, tt.getObject, nil, nil, nil)

			routerClients := router.Clients{
				S3Client: s3Client,
			}

			routerUUIDGen := router.UUIDGenerator{
				UUIDGen: uuidgenerator.NewUuidGeneratorDummy(nil, tt.isValidUUID),
			}

			r := router.NewRouter(config.Config{
				UserAuth: givenUsername,
				PwdAuth:  givenUserPwd,
			}, &routerClients, &routerUUIDGen, &router.DAOs{})

			w := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, tt.giveRequest, nil)
			if tt.giveWithAuth {
				req.SetBasicAuth(givenUsername, givenUserPwd)
			}

			r.ServeHTTP(w, req)
			require.Equal(t, tt.expectedHTTPCode, w.Code)
		})
	}

}
