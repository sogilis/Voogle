package controllers_test

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	. "github.com/Sogilis/Voogle/services/api/controllers"
)

func TestVideosListHandler(t *testing.T) {
	// Given
	assert.Nil(t, os.RemoveAll("./videos"))
	assert.Nil(t, os.Mkdir("./videos", os.ModePerm))
	assert.Nil(t, os.Mkdir("./videos/video1", os.ModePerm))
	assert.Nil(t, os.Mkdir("./videos/video2", os.ModePerm))

	allVideosExpected := AllVideos{Status: "Success", Data: []VideoInfo{{Id: "video1", Title: "video1"}, {Id: "video2", Title: "video2"}}}
	w := httptest.NewRecorder()

	// When
	r := mux.NewRouter()
	r.PathPrefix("/api/v1/videos").Handler(VideosListHandler{}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/v1/videos", nil)
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, 200, w.Code)

	gotAllVideos := AllVideos{}
	assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &gotAllVideos))

	assert.True(t, reflect.DeepEqual(allVideosExpected, gotAllVideos))
	assert.Nil(t, os.RemoveAll("./videos"))
}
