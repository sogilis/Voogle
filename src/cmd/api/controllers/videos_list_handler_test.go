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
	"github.com/stretchr/testify/require"

	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	. "github.com/Sogilis/Voogle/src/cmd/api/controllers"
	"github.com/Sogilis/Voogle/src/cmd/api/models"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
)

func TestVideosListHandler(t *testing.T) { //nolint:cyclop
	// Given
	videosExpected := AllVideos{
		Status: "Success",
		Data: []VideoInfo{
			{Id: uuid.NewString(), Title: "video1"},
			{Id: uuid.NewString(), Title: "video2"},
		},
	}
	w := httptest.NewRecorder()

	testUsername := "dev"
	testUsePwd := "test"

	// Mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	routerClients := router.Clients{
		MariadbClient: db,
	}

	routerUUIDGen := router.UUIDGenerator{
		UUIDGen: uuidgenerator.NewUuidGeneratorDummy(nil, nil),
	}

	// Queries
	getVideos := regexp.QuoteMeta("SELECT * FROM videos v")

	// Tables
	videosColumns := []string{"id", "title", "video_status", "uploaded_at", "created_at", "updated_at"}
	videosRows := sqlmock.NewRows(videosColumns)

	t1 := time.Now()
	videosRows.AddRow(videosExpected.Data[0].Id, videosExpected.Data[0].Title, models.COMPLETE, nil, t1, nil)
	videosRows.AddRow(videosExpected.Data[1].Id, videosExpected.Data[1].Title, models.COMPLETE, nil, t1, nil)

	mock.ExpectPrepare(getVideos)
	mock.ExpectQuery(getVideos).WillReturnRows(videosRows)

	// When
	r := router.NewRouter(config.Config{
		UserAuth: testUsername,
		PwdAuth:  testUsePwd,
	}, &routerClients, &routerUUIDGen)

	req := httptest.NewRequest("GET", "/api/v1/videos/list", nil)
	req.SetBasicAuth(testUsername, testUsePwd)
	r.ServeHTTP(w, req)

	// Then
	require.Equal(t, 200, w.Code)

	gotAllVideos := AllVideos{}
	require.Nil(t, json.Unmarshal(w.Body.Bytes(), &gotAllVideos))

	require.True(t, reflect.DeepEqual(videosExpected, gotAllVideos))

	// we make sure that all expectations were met
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
