package controllers_test

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

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
	}{
		{name: "GET video stream master", request: "/api/v1/videos/" + validVideoID + "/streams/master.m3u8", withAuth: true, expectedHTTPCode: 200},
		{name: "GET fails to video stream master", request: "/api/v1/videos/" + "invalidID" + "/streams/master.m3u8", withAuth: true, expectedHTTPCode: 404},
		{name: "GET video sub part", request: "/api/v1/videos/" + validVideoID + "/streams/" + validQuality + "/" + validSubPart, withAuth: true, expectedHTTPCode: 200},
		{name: "GET fails with wrong quality", request: "/api/v1/videos/" + validVideoID + "/streams/" + "v1" + "/" + validSubPart, withAuth: true, expectedHTTPCode: 404},
		{name: "GET fails with wrong subpart name", request: "/api/v1/videos/" + validVideoID + "/streams/" + validQuality + "/" + "invalidSubPart", withAuth: true, expectedHTTPCode: 404},
		{name: "GET fails with no auth", request: "/api/v1/videos/" + validVideoID + "/streams/" + validQuality + "/" + "invalidSubPart", withAuth: false, expectedHTTPCode: 401},
	}

	r := NewRouter(config.Config{
		UserAuth: givenUsername,
		PwdAuth:  givenUserPwd,
	})

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			buildDummyDirectory(t, validVideoID, validQuality, validSubPart)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", tt.request, nil)
			if tt.withAuth {
				req.SetBasicAuth(givenUsername, givenUserPwd)
			}

			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedHTTPCode, w.Code)

			removeDummyDirectory(t)
		})
	}

}

/// Create a dummy substructure of a HLS video
func buildDummyDirectory(t *testing.T, videoID, quality, part string) {
	assert.Nil(t, os.RemoveAll("./videos"))
	assert.Nil(t, os.MkdirAll("./videos/"+videoID+"/"+quality, os.ModePerm))

	f, err := os.Create("./videos/" + videoID + "/master.m3u8")
	assert.Nil(t, err)
	assert.Nil(t, f.Close())
	f, err = os.Create("./videos/" + videoID + "/" + quality + "/" + part)
	assert.Nil(t, err)
	assert.Nil(t, f.Close())
}

func removeDummyDirectory(t *testing.T) {
	assert.Nil(t, os.RemoveAll("./videos"))
}
