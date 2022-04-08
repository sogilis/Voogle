package controllers_test

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	. "github.com/Sogilis/Voogle/src/cmd/api/controllers"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
	. "github.com/Sogilis/Voogle/src/cmd/api/router"
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

	t1 := time.Now()
	rows := sqlmock.NewRows([]string{"id", "title", "v_status", "uploaded_at", "created_at", "updated_at"}).
		AddRow(allVideosExpected.Data[0].Id, allVideosExpected.Data[0].Title, int(contracts.Video_VIDEO_STATUS_ENCODING), nil, t1, t1).
		AddRow(allVideosExpected.Data[1].Id, allVideosExpected.Data[1].Title, int(contracts.Video_VIDEO_STATUS_ENCODING), nil, t1, t1)

	query := regexp.QuoteMeta("SELECT * FROM videos v")

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
