package controllers_test

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	. "github.com/Sogilis/Voogle/src/cmd/api/controllers"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
	. "github.com/Sogilis/Voogle/src/cmd/api/router"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"
)

func TestVideosListHandler(t *testing.T) {
	// Given
	allVideosExpected := AllVideos{Status: "Success", Data: []VideoInfo{{Id: uuid.NewString(), Title: "video1"}, {Id: uuid.NewString(), Title: "video2"}}}
	w := httptest.NewRecorder()

	testUsername := "dev"
	testUsePwd := "test"

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "public_id", "title", "state_name", "last_update"}).
		AddRow(uuid.NewString(), allVideosExpected.Data[0].Id, allVideosExpected.Data[0].Title, "UPLOADING", time.Now()).
		AddRow(uuid.NewString(), allVideosExpected.Data[1].Id, allVideosExpected.Data[1].Title, "UPLOADING", time.Now())

	query := `SELECT v.id, public_id, title, state_name, last_update
			  FROM videos v
			  INNER JOIN video_state vs ON v.v_state = vs.id;`

	mock.ExpectQuery(query).WillReturnRows(rows)

	routerClients := router.Clients{
		MariadbClient: db,
	}

	uuidGen := uuidgenerator.NewUuidGeneratorDummy(nil)

	routerUUIDGen := UUIDGenerator{
		UUIDGen: uuidGen,
	}

	// When
	r := NewRouter(config.Config{
		UserAuth: testUsername,
		PwdAuth:  testUsePwd,
	}, &routerClients, &routerUUIDGen)

	req := httptest.NewRequest("GET", "/api/v1/videos/list", nil)
	req.SetBasicAuth(testUsername, testUsePwd)
	r.ServeHTTP(w, req)

	// Then
	assert.Equal(t, 200, w.Code)

	gotAllVideos := AllVideos{}
	assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &gotAllVideos))

	assert.True(t, reflect.DeepEqual(allVideosExpected, gotAllVideos))
}
