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
	assert.Nil(t, os.RemoveAll("./videos"))
	assert.Nil(t, os.Mkdir("./videos", os.ModePerm))
	assert.Nil(t, os.Mkdir("./videos/video1", os.ModePerm))
	f, err := os.Create("./videos/video1/master.m3u8")
	assert.Nil(t, err)
	err = f.Close()
	assert.Nil(t, err)

	w := httptest.NewRecorder()

	// When
	r := mux.NewRouter()
	r.PathPrefix("/api/v1/videos/{id}/streams/master.m3u8").Handler(VideoGetMasterHandler{}).Methods("GET")

	req := httptest.NewRequest("GET", "/api/v1/videos/video1/streams/master.m3u8", nil)
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, 200, w.Code)
	assert.Nil(t, os.RemoveAll("./videos"))
}
