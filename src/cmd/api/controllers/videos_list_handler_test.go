package controllers_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	"github.com/Sogilis/Voogle/src/cmd/api/db/dao_test"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
)

func TestVideosList(t *testing.T) { //nolint:cyclop

	//Initialize and set default parameters
	givenUsername := "dev"
	givenPassword := "test"
	validVideoId := "1508e7d5-5bc6-4a50-9176-ab0371aa65fe"
	t1 := time.Now()
	UUIDValidFunc := func(u string) bool { _, err := uuid.Parse(u); return err == nil }

	cases := []struct {
		name             string
		authIsGiven      bool
		databaseHasError bool
		requestIsWrong   bool
		videoAttribute   string
		ascending        string
		page             string
		limit            string
		expectedHTTPCode int
	}{
		{
			name:             "GET videos list",
			authIsGiven:      true,
			databaseHasError: false,
			videoAttribute:   "title",
			ascending:        "true",
			page:             "1",
			limit:            "10",
			expectedHTTPCode: 200,
		},
		{
			name:             "GET fails with missing attribute",
			authIsGiven:      true,
			databaseHasError: true,
			requestIsWrong:   true,
			videoAttribute:   "invalid",
			ascending:        "true",
			page:             "1",
			limit:            "10",
			expectedHTTPCode: 400,
		},
		{
			name:             "GET fails with missing order",
			authIsGiven:      true,
			databaseHasError: true,
			requestIsWrong:   true,
			videoAttribute:   "title",
			ascending:        "invalid",
			page:             "1",
			limit:            "10",
			expectedHTTPCode: 400,
		},
		{
			name:             "GET fails with missing page number",
			authIsGiven:      true,
			databaseHasError: true,
			requestIsWrong:   true,
			videoAttribute:   "title",
			ascending:        "true",
			page:             "invalid",
			limit:            "10",
			expectedHTTPCode: 400,
		},
		{
			name:             "GET fails with missing limit number",
			authIsGiven:      true,
			databaseHasError: true,
			requestIsWrong:   true,
			videoAttribute:   "title",
			ascending:        "true",
			page:             "1",
			limit:            "invalid",
			expectedHTTPCode: 400,
		},
		{
			name:             "GET fails with authentification error",
			authIsGiven:      false,
			databaseHasError: false,
			videoAttribute:   "title",
			ascending:        "true",
			page:             "1",
			limit:            "10",
			expectedHTTPCode: 401,
		},
		{
			name:             "GET fails with database error",
			authIsGiven:      true,
			databaseHasError: true,
			videoAttribute:   "title",
			ascending:        "true",
			page:             "1",
			limit:            "10",
			expectedHTTPCode: 500,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			// Mock database
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			routerUUIDGen := router.UUIDGenerator{
				UUIDGen: uuidgenerator.NewUuidGeneratorDummy(nil, UUIDValidFunc),
			}

			// Init videoDAO
			dao_test.ExpectVideosDAOCreation(mock)

			//Create request
			givenRequest := fmt.Sprintf("/api/v1/videos/list/%v/%v/%v/%v", tt.videoAttribute, tt.ascending, tt.page, tt.limit)

			if !tt.authIsGiven || tt.requestIsWrong {
				// This case will stop before modifying the database

			} else {

				direction := "DESC"
				if tt.ascending == "true" {
					direction = "ASC"
				}
				pagenum, _ := strconv.Atoi(tt.page)
				limitnum, _ := strconv.Atoi(tt.limit)
				// Queries
				getVideoListQuery := regexp.QuoteMeta(fmt.Sprintf("SELECT * FROM videos ORDER BY %v %v LIMIT %d,%d", tt.videoAttribute, direction, (pagenum-1)*limitnum, limitnum))
				getVideoTotal := regexp.QuoteMeta(dao.VideosRequests[dao.GetTotalVideos])

				// Tables
				videosColumns := []string{"id", "title", "video_status", "uploaded_at", "created_at", "updated_at", "source_path"}
				videosRows := sqlmock.NewRows(videosColumns)

				if tt.databaseHasError {
					// mock.ExpectQuery(getVideoListQuery).WillReturnError(fmt.Errorf("Server Error"))
					mock.ExpectQuery(getVideoListQuery).WillReturnError(fmt.Errorf("Server Error"))
				} else {
					sourcePathVideo := validVideoId + "/" + "source.mp4"
					videosRows.AddRow(validVideoId, "title", contracts.Video_VIDEO_STATUS_ENCODING, t1, t1, nil, sourcePathVideo)
					mock.ExpectQuery(getVideoListQuery).WillReturnRows(videosRows)
					mock.ExpectQuery(getVideoTotal).WillReturnRows(sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(1))
				}
			}

			// When
			VideosDAO, err := dao.CreateVideosDAO(context.Background(), db)
			require.NoError(t, err)
			routerDAO := router.DAOs{
				VideosDAO: *VideosDAO,
			}

			r := router.NewRouter(config.Config{
				UserAuth: givenUsername,
				PwdAuth:  givenPassword,
			}, &router.Clients{}, &routerUUIDGen, &routerDAO)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", givenRequest, nil)
			if tt.authIsGiven {
				req.SetBasicAuth(givenUsername, givenPassword)
			}

			r.ServeHTTP(w, req)
			require.Equal(t, tt.expectedHTTPCode, w.Code)

			// we make sure that all expectations were met
			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}
