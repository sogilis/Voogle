package controllers_test

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Sogilis/Voogle/services/api/clients"
	"github.com/Sogilis/Voogle/services/api/config"
	. "github.com/Sogilis/Voogle/services/api/controllers"
	"github.com/Sogilis/Voogle/services/api/router"
	. "github.com/Sogilis/Voogle/services/api/router"
)

func TestVideosListHandler(t *testing.T) {
	// Given
	allVideosExpected := AllVideos{Status: "Success", Data: []VideoInfo{{Id: "video1", Title: "video1"}, {Id: "video2", Title: "video2"}}}
	w := httptest.NewRecorder()

	testUsername := "dev"
	testUsePwd := "test"

	s3Client := clients.NewS3ClientDummy(func() ([]string, error) {
		return []string{"video1", "video2"}, nil
	}, nil)

	routerClients := router.Clients{
		S3Client: s3Client,
	}

	// When
	r := NewRouter(config.Config{
		UserAuth: testUsername,
		PwdAuth:  testUsePwd,
	}, &routerClients)

	req := httptest.NewRequest("GET", "/api/v1/videos/list", nil)
	req.SetBasicAuth(testUsername, testUsePwd)
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, 200, w.Code)

	gotAllVideos := AllVideos{}
	assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &gotAllVideos))

	assert.True(t, reflect.DeepEqual(allVideosExpected, gotAllVideos))
}
