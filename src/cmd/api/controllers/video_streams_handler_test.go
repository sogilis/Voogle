package controllers_test

import (
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	. "github.com/Sogilis/Voogle/src/cmd/api/router"
)

func TestVideoServe(t *testing.T) {
	givenUsername := "dev"
	givenUserPwd := "test"

	validVideoID := "video1"
	validQuality := "v0"
	validSubPart := "part1.ts"

	cases := []struct {
		name             string
		giveRequest      string
		giveWithAuth     bool
		expectedHTTPCode int
		getObjectID      func(string) (io.Reader, error)
	}{
		{name: "GET video stream master", giveRequest: "/api/v1/videos/" + validVideoID + "/streams/master.m3u8", giveWithAuth: true, expectedHTTPCode: 200, getObjectID: func(s string) (io.Reader, error) { return strings.NewReader(""), nil }},
		{name: "GET fails to video stream master", giveRequest: "/api/v1/videos/" + "invalidID" + "/streams/master.m3u8", giveWithAuth: true, expectedHTTPCode: 404, getObjectID: func(s string) (io.Reader, error) { return nil, errors.New("Not found") }},
		{name: "GET video sub part", giveRequest: "/api/v1/videos/" + validVideoID + "/streams/" + validQuality + "/" + validSubPart, giveWithAuth: true, expectedHTTPCode: 200, getObjectID: func(s string) (io.Reader, error) { return strings.NewReader(""), nil }},
		{name: "GET fails with wrong quality", giveRequest: "/api/v1/videos/" + validVideoID + "/streams/" + "v1" + "/" + validSubPart, giveWithAuth: true, expectedHTTPCode: 404, getObjectID: func(s string) (io.Reader, error) { return nil, errors.New("Not found") }},
		{name: "GET fails with wrong subpart name", giveRequest: "/api/v1/videos/" + validVideoID + "/streams/" + validQuality + "/" + "invalidSubPart", giveWithAuth: true, expectedHTTPCode: 404, getObjectID: func(s string) (io.Reader, error) { return nil, errors.New("Not found") }},
		{name: "GET fails with no auth", giveRequest: "/api/v1/videos/" + validVideoID + "/streams/" + validQuality + "/" + "invalidSubPart", giveWithAuth: false, expectedHTTPCode: 401, getObjectID: nil},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			s3Client := clients.NewS3ClientDummy(nil, tt.getObjectID, nil, nil)

			routerClients := Clients{
				S3Client: s3Client,
			}

			routerUUIDGen := UUIDGenerator{
				UUIDGen: uuidgenerator.NewUuidGeneratorDummy(nil),
			}

			r := NewRouter(config.Config{
				UserAuth: givenUsername,
				PwdAuth:  givenUserPwd,
			}, &routerClients, &routerUUIDGen)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", tt.giveRequest, nil)
			if tt.giveWithAuth {
				req.SetBasicAuth(givenUsername, givenUserPwd)
			}

			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedHTTPCode, w.Code)

		})

	}

}
