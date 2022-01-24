package controllers_test

import (
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Sogilis/Voogle/services/api/clients"
	"github.com/Sogilis/Voogle/services/api/config"
	. "github.com/Sogilis/Voogle/services/api/router"
)

func TestVideoServe(t *testing.T) {
	givenUsername := "dev"
	givenUserPwd := "test"

	validVideoID := "video1"
	validQuality := "v0"
	validSubPart := "part1.ts"

	cases := []struct {
		name             string
		request          string
		withAuth         bool
		expectedHTTPCode int
		getObjectID      func(string) (io.Reader, error)
	}{
		{name: "GET video stream master", request: "/api/v1/videos/" + validVideoID + "/streams/master.m3u8", withAuth: true, expectedHTTPCode: 200, getObjectID: func(s string) (io.Reader, error) { return strings.NewReader(""), nil }},
		{name: "GET fails to video stream master", request: "/api/v1/videos/" + "invalidID" + "/streams/master.m3u8", withAuth: true, expectedHTTPCode: 404, getObjectID: func(s string) (io.Reader, error) { return nil, errors.New("Not found") }},
		{name: "GET video sub part", request: "/api/v1/videos/" + validVideoID + "/streams/" + validQuality + "/" + validSubPart, withAuth: true, expectedHTTPCode: 200, getObjectID: func(s string) (io.Reader, error) { return strings.NewReader(""), nil }},
		{name: "GET fails with wrong quality", request: "/api/v1/videos/" + validVideoID + "/streams/" + "v1" + "/" + validSubPart, withAuth: true, expectedHTTPCode: 404, getObjectID: func(s string) (io.Reader, error) { return nil, errors.New("Not found") }},
		{name: "GET fails with wrong subpart name", request: "/api/v1/videos/" + validVideoID + "/streams/" + validQuality + "/" + "invalidSubPart", withAuth: true, expectedHTTPCode: 404, getObjectID: func(s string) (io.Reader, error) { return nil, errors.New("Not found") }},
		{name: "GET fails with no auth", request: "/api/v1/videos/" + validVideoID + "/streams/" + validQuality + "/" + "invalidSubPart", withAuth: false, expectedHTTPCode: 401, getObjectID: nil},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			s3Client := clients.NewS3ClientDummy(nil, tt.getObjectID)

			routerClients := Clients{
				S3Client: s3Client,
			}

			r := NewRouter(config.Config{
				UserAuth: givenUsername,
				PwdAuth:  givenUserPwd,
			}, &routerClients)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", tt.request, nil)
			if tt.withAuth {
				req.SetBasicAuth(givenUsername, givenUserPwd)
			}

			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedHTTPCode, w.Code)

		})

	}

}
