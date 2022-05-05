package controllers_test

import (
	"fmt"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestVideosListHandler(t *testing.T) { //nolint:cyclop

	//Initialize and set default parameters
	givenUsername := "dev"
	givenPassword := "test"
	validVideoId := "1508e7d5-5bc6-4a50-9176-ab0371aa65fe"
	videoTitle := "title"
	ascending := true
	page := 1
	limit := 10
	t1 := time.Now()
	UUIDValidFunc := func(u string) bool { _, err := uuid.Parse(u); return err == nil }

	cases := []struct {
		name             string
		authIsGiven      bool
		databaseHasError bool
		givenQuery       string
		expectedHTTPCode int
	}{
		{
			name:             "GET videos list",
			authIsGiven:      true,
			databaseHasError: false,
			givenQuery:       fmt.Sprintf("/api/v1/videos/list/%v/%v/%v/%v", videoTitle, ascending, page, limit),
			expectedHTTPCode: 200,
		},
		{
			name:             "GET fails with authentification error",
			authIsGiven:      true,
			databaseHasError: false,
			givenQuery:       fmt.Sprintf("/api/v1/videos/list/%v/%v/%v/%v", videoTitle, ascending, page, limit),
			expectedHTTPCode: 401,
		},
		{
			name:             "GET fails with database error",
			authIsGiven:      true,
			databaseHasError: true,
			givenQuery:       fmt.Sprintf("/api/v1/videos/list/%v/%v/%v/%v", videoTitle, ascending, page, limit),
			expectedHTTPCode: 500,
		},
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
				UUIDGen: uuidgenerator.NewUuidGeneratorDummy(nil, UUIDValidFunc),
			}

			if !tt.authIsGiven {
				// This case will stop before modifying the database

			} else {

				// Queries
				getVideoListQuery := regexp.QuoteMeta("SELECT * FROM videos ORDER BY ? ? LIMIT ?,?")

				// Tables
				videosColumns := []string{"id", "title", "video_status", "uploaded_at", "created_at", "updated_at"}
				videosRows := sqlmock.NewRows(videosColumns)

				if tt.databaseHasError {
					mock.ExpectQuery(getVideoListQuery).WithArgs("title", "ASC", "0", "10").WillReturnError(fmt.Errorf("Server Error"))
				} else {
					videosRows.AddRow(validVideoId, videoTitle, contracts.Video_VIDEO_STATUS_ENCODING, t1, t1, nil)
					mock.ExpectQuery(getVideoListQuery).WithArgs("title", "ASC", "0", "10").WillReturnRows(videosRows)
				}
			}

			r := router.NewRouter(config.Config{
				UserAuth: givenUsername,
				PwdAuth:  givenPassword,
			}, &routerClients, &routerUUIDGen)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", tt.givenQuery, nil)
			if tt.authIsGiven {
				req.SetBasicAuth(givenUsername, givenPassword)
			}

			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedHTTPCode, w.Code)

			// we make sure that all expectations were met
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
