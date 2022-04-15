package controllers_test

import (
	"fmt"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
	. "github.com/Sogilis/Voogle/src/cmd/api/router"
)

func TestVideoStatus(t *testing.T) {
	givenUsername := "dev"
	givenUserPwd := "test"

	validVideoID := "video1"
	videoTitle := "title"
	t1 := time.Now()

	cases := []struct {
		name             string
		giveRequest      string
		giveWithAuth     bool
		expectedHTTPCode int
	}{
		{
			name:             "GET video status",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/status",
			giveWithAuth:     true,
			expectedHTTPCode: 200},
		{
			name:             "GET fails with invalid video ID",
			giveRequest:      "/api/v1/videos/" + "invalidID" + "/status",
			giveWithAuth:     true,
			expectedHTTPCode: 400},

		{
			name:             "GET fails with no auth",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/status",
			giveWithAuth:     false,
			expectedHTTPCode: 401},
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

			routerUUIDGen := UUIDGenerator{
				UUIDGen: uuidgenerator.NewUuidGeneratorDummy(nil),
			}

			if !tt.giveWithAuth {
				// All these cases will stop before modifying the database : Nothing to do

			} else {
				// Queries
				getVideoFromIdQuery := regexp.QuoteMeta("SELECT * FROM videos v WHERE v.id = ?")

				// Tables
				videosColumns := []string{"id", "title", "video_status", "uploaded_at", "created_at", "updated_at"}
				videosRows := sqlmock.NewRows(videosColumns)

				if tt.giveRequest == "/api/v1/videos/"+"invalidID"+"/status" {
					mock.ExpectQuery(getVideoFromIdQuery).WillReturnError(fmt.Errorf("unknow invalid video ID"))
				} else {
					videosRows.AddRow(validVideoID, videoTitle, contracts.Video_VIDEO_STATUS_ENCODING, nil, t1, nil)
					mock.ExpectQuery(getVideoFromIdQuery).WillReturnRows(videosRows)
				}
			}

			r := NewRouter(config.Config{
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
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

		})

	}

}
