package controllers_test

import (
	"errors"
	"io"
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

func TestVideoServe(t *testing.T) { //nolint:cyclop
	givenUsername := "dev"
	givenUserPwd := "test"

	validVideoID := "1508e7d5-5bc6-4a50-9176-ab0371aa65fe"
	invalidVideoID := "invalidvideoid"
	unknownVideoID := "0000a0a0-0aa0-0a00-0000-aa0000aa00aa"
	validQuality := "v0"
	validSubPart := "part1.ts"
	UUIDValidFunc := func(u string) bool { _, err := uuid.Parse(u); return err == nil }

	cases := []struct {
		name             string
		giveRequest      string
		giveWithAuth     bool
		expectedHTTPCode int
		getObjectID      func(string) (io.Reader, error)
		isValidUUID      func(string) bool
	}{
		{
			name:             "GET video stream master",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/streams/master.m3u8",
			giveWithAuth:     true,
			expectedHTTPCode: 200,
			getObjectID:      func(s string) (io.Reader, error) { return strings.NewReader(""), nil },
			isValidUUID:      UUIDValidFunc},
		{
			name:             "GET fails to video stream master with invalid id",
			giveRequest:      "/api/v1/videos/" + invalidVideoID + "/streams/master.m3u8",
			giveWithAuth:     true,
			expectedHTTPCode: 400,
			getObjectID:      func(s string) (io.Reader, error) { return strings.NewReader(""), nil },
			isValidUUID:      UUIDValidFunc},
		{
			name:             "GET fails to video stream master with unknown id",
			giveRequest:      "/api/v1/videos/" + unknownVideoID + "/streams/master.m3u8",
			giveWithAuth:     true,
			expectedHTTPCode: 404,
			getObjectID:      func(s string) (io.Reader, error) { return nil, errors.New("Not found") },
			isValidUUID:      UUIDValidFunc},
		{
			name:             "GET video sub part",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/streams/" + validQuality + "/" + validSubPart,
			giveWithAuth:     true,
			expectedHTTPCode: 200,
			getObjectID:      func(s string) (io.Reader, error) { return strings.NewReader(""), nil },
			isValidUUID:      UUIDValidFunc},
		{
			name:             "GET fails with wrong quality",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/streams/" + "v1" + "/" + validSubPart,
			giveWithAuth:     true,
			expectedHTTPCode: 404,
			getObjectID:      func(s string) (io.Reader, error) { return nil, errors.New("Not found") },
			isValidUUID:      UUIDValidFunc},
		{
			name:             "GET fails with wrong subpart name",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/streams/" + validQuality + "/" + "invalidSubPart",
			giveWithAuth:     true,
			expectedHTTPCode: 404,
			getObjectID:      func(s string) (io.Reader, error) { return nil, errors.New("Not found") },
			isValidUUID:      UUIDValidFunc},
		{
			name:             "GET fails with invalid video ID",
			giveRequest:      "/api/v1/videos/" + invalidVideoID + "/streams/" + validQuality + "/" + validSubPart,
			giveWithAuth:     true,
			expectedHTTPCode: 400,
			getObjectID:      func(s string) (io.Reader, error) { return nil, errors.New("Not found") },
			isValidUUID:      UUIDValidFunc},
		{
			name:             "GET fails with unknow video ID",
			giveRequest:      "/api/v1/videos/" + unknownVideoID + "/streams/" + validQuality + "/" + validSubPart,
			giveWithAuth:     true,
			expectedHTTPCode: 404,
			getObjectID:      func(s string) (io.Reader, error) { return nil, errors.New("Not found") },
			isValidUUID:      UUIDValidFunc},
		{
			name:             "GET fails with no auth",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/streams/" + validQuality + "/" + "invalidSubPart",
			giveWithAuth:     false,
			expectedHTTPCode: 401,
			getObjectID:      nil},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			s3Client := clients.NewS3ClientDummy(nil, tt.getObjectID, nil, nil, nil)

			routerClients := router.Clients{
				S3Client: s3Client,
			}

			routerUUIDGen := router.UUIDGenerator{
				UUIDGen: uuidgenerator.NewUuidGeneratorDummy(nil, tt.isValidUUID),
			}

			r := router.NewRouter(config.Config{
				UserAuth: givenUsername,
				PwdAuth:  givenUserPwd,
			}, &routerClients, &routerUUIDGen)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", tt.giveRequest, nil)
			if tt.giveWithAuth {
				req.SetBasicAuth(givenUsername, givenUserPwd)
			}

			r.ServeHTTP(w, req)
			require.Equal(t, tt.expectedHTTPCode, w.Code)

		})

	}

}
