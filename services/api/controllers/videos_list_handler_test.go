package controllers_test

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Sogilis/Voogle/services/api/config"
	. "github.com/Sogilis/Voogle/services/api/controllers"
	. "github.com/Sogilis/Voogle/services/api/router"
)

func TestVideosListHandler(t *testing.T) {
	// Given
	assert.Nil(t, os.RemoveAll("./videos"))
	assert.Nil(t, os.Mkdir("./videos", os.ModePerm))
	assert.Nil(t, os.Mkdir("./videos/video1", os.ModePerm))
	assert.Nil(t, os.Mkdir("./videos/video2", os.ModePerm))

	allVideosExpected := AllVideos{Status: "Success", Data: []VideoInfo{{Id: "video1", Title: "video1"}, {Id: "video2", Title: "video2"}}}
	w := httptest.NewRecorder()

	testUsername := "dev"
	testUsePwd := "test"

	// When
	r := NewRouter(config.Config{
		UserAuth: testUsername,
		PwdAuth:  testUsePwd,
	})

	req := httptest.NewRequest("GET", "/api/v1/videos/list", nil)
	req.SetBasicAuth(testUsername, testUsePwd)
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, 200, w.Code)

	gotAllVideos := AllVideos{}
	assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &gotAllVideos))

	assert.True(t, reflect.DeepEqual(allVideosExpected, gotAllVideos))
	assert.Nil(t, os.RemoveAll("./videos"))
}
