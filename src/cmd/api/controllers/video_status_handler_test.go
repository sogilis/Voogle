package controllers_test

import (
	"fmt"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
)

func TestVideoStatus(t *testing.T) { //nolint:cyclop
	givenUsername := "dev"
	givenUserPwd := "test"

	validVideoID := "1508e7d5-5bc6-4a50-9176-ab0371aa65fe"
	invalidVideoID := "invalidvideoid"
	unknownVideoID := "0000a0a0-0aa0-0a00-0000-aa0000aa00aa"
	UUIDValidFunc := func(u string) bool { _, err := uuid.Parse(u); return err == nil }
	videoTitle := "title"
	t1 := time.Now()

	cases := []struct {
		name             string
		giveRequest      string
		giveWithAuth     bool
		giveDatabaseErr  bool
		expectedHTTPCode int
		isValidUUID      func(string) bool
	}{
		{
			name:             "GET video status",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/status",
			giveWithAuth:     true,
			expectedHTTPCode: 200,
			isValidUUID:      UUIDValidFunc},
		{
			name:             "GET fails with invalid video ID",
			giveRequest:      "/api/v1/videos/" + invalidVideoID + "/status",
			giveWithAuth:     true,
			expectedHTTPCode: 400,
			isValidUUID:      UUIDValidFunc},
		{
			name:             "GET fails with unknown video ID",
			giveRequest:      "/api/v1/videos/" + unknownVideoID + "/status",
			giveWithAuth:     true,
			expectedHTTPCode: 404,
			isValidUUID:      UUIDValidFunc},
		{
			name:             "GET fails with database error",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/status",
			giveWithAuth:     true,
			giveDatabaseErr:  true,
			expectedHTTPCode: 500,
			isValidUUID:      UUIDValidFunc},

		{
			name:             "GET fails with no auth",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/status",
			giveWithAuth:     false,
			expectedHTTPCode: 401,
			isValidUUID:      UUIDValidFunc},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			// Mock database
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			routerClients := router.Clients{
				MariadbClient: db,
			}

			routerUUIDGen := router.UUIDGenerator{
				UUIDGen: uuidgenerator.NewUuidGeneratorDummy(nil, tt.isValidUUID),
			}

			if !tt.giveWithAuth || tt.giveRequest == "/api/v1/videos/"+invalidVideoID+"/status" {
				// All these cases will stop before modifying the database : Nothing to do

			} else {
				// Queries
				getVideoFromIdQuery := regexp.QuoteMeta("SELECT * FROM videos v WHERE v.id = ?")

				// Tables
				videosColumns := []string{"id", "title", "video_status", "uploaded_at", "created_at", "updated_at"}
				videosRows := sqlmock.NewRows(videosColumns)

				if tt.giveDatabaseErr {
					mock.ExpectQuery(getVideoFromIdQuery).WillReturnError(fmt.Errorf("unknow invalid video ID"))

				} else if tt.giveRequest == "/api/v1/videos/"+unknownVideoID+"/status" {
					mock.ExpectQuery(getVideoFromIdQuery).WillReturnRows(videosRows)

				} else {
					videosRows.AddRow(validVideoID, videoTitle, contracts.Video_VIDEO_STATUS_ENCODING, nil, t1, nil)
					mock.ExpectQuery(getVideoFromIdQuery).WillReturnRows(videosRows)
				}
			}

			r := router.NewRouter(config.Config{
				UserAuth: givenUsername,
				PwdAuth:  givenUserPwd,
			}, &routerClients, &routerUUIDGen)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", tt.giveRequest, nil)
			if tt.giveWithAuth {
				req.SetBasicAuth(givenUsername, givenUserPwd)
			}

			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedHTTPCode, w.Code)

			// we make sure that all expectations were met
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})

	}

}
