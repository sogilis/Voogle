package controllers_test

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	. "github.com/Sogilis/Voogle/services/api/controllers"
)

func TestVideosStreamsMasterHandler(t *testing.T) {
	// Given
	videoID := "video1"
	quality := "v0"
	part := "part1.ts"

	buildDummyDirectory(t, videoID, quality, part)

	w := httptest.NewRecorder()

	// When
	r := mux.NewRouter()
	r.PathPrefix("/api/v1/videos/{id}/streams/master.m3u8").Handler(VideoGetMasterHandler{}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/v1/videos/"+videoID+"/streams/master.m3u8", nil)
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, 200, w.Code)
	removeDummyDirectory(t)
}

func TestGetVideoMasterFails(t *testing.T) {
	// Given
	videoID := "video1"
	quality := "v0"
	part := "part1.ts"

	buildDummyDirectory(t, videoID, quality, part)

	w := httptest.NewRecorder()

	// When
	r := mux.NewRouter()
	r.PathPrefix("/api/v1/videos/{id}/streams/master.m3u8").Handler(VideoGetMasterHandler{}).Methods("GET")

	wantVideoID := "none"
	req := httptest.NewRequest("GET", "/api/v1/videos/"+wantVideoID+"/streams/master.m3u8", nil)
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, 404, w.Code)
	removeDummyDirectory(t)
}

func TestGetVideoSubPartSucceed(t *testing.T) {
	// Given
	videoID := "video1"
	quality := "v0"
	part := "part1.ts"

	buildDummyDirectory(t, videoID, quality, part)

	w := httptest.NewRecorder()

	// When
	r := mux.NewRouter()
	r.PathPrefix("/api/v1/videos/{id}/streams/{quality}/{filename}").Handler(VideoGetSubPartHandler{}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/v1/videos/"+videoID+"/streams/"+quality+"/"+part, nil)
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, 200, w.Code)
	removeDummyDirectory(t)
}

func TestGetVideoSubPartFails(t *testing.T) {
	// Given
	videoID := "video1"
	quality := "v0"
	part := "part1.ts"

	buildDummyDirectory(t, videoID, quality, part)

	w := httptest.NewRecorder()

	// When
	r := mux.NewRouter()
	r.PathPrefix("/api/v1/videos/{id}/streams/{quality}/{filename}").Handler(VideoGetSubPartHandler{}).Methods("GET")

	wantQuality := "None"
	req := httptest.NewRequest("GET", "/api/v1/videos/"+videoID+"/streams/"+wantQuality+"/"+part, nil)
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, 404, w.Code)
	removeDummyDirectory(t)
}

/// Create a dummy substructure of a HLS video
func buildDummyDirectory(t *testing.T, videoID, quality, part string) {
	assert.Nil(t, os.RemoveAll("./videos"))
	os.MkdirAll("./videos/"+videoID+"/"+quality, os.ModePerm)

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
