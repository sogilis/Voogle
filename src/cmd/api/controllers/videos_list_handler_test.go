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

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	"github.com/Sogilis/Voogle/src/cmd/api/db/dao_test"
	"github.com/Sogilis/Voogle/src/cmd/api/models"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
	"github.com/Sogilis/Voogle/src/pkg/clients"
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
		status           models.VideoStatus
		expectedHTTPCode int
		isValidUUID      func(string) bool
	}{
		{
			name:             "GET complete videos list",
			authIsGiven:      true,
			databaseHasError: false,
			videoAttribute:   "title",
			ascending:        "true",
			page:             "1",
			limit:            "10",
			status:           models.COMPLETE,
			expectedHTTPCode: 200,
			isValidUUID:      UUIDValidFunc,
		},
		{
			name:             "GET archived videos list",
			authIsGiven:      true,
			databaseHasError: false,
			videoAttribute:   "title",
			ascending:        "true",
			page:             "1",
			limit:            "10",
			status:           models.ARCHIVE,
			expectedHTTPCode: 200,
			isValidUUID:      UUIDValidFunc,
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
			status:           models.COMPLETE,
			expectedHTTPCode: 400,
			isValidUUID:      UUIDValidFunc,
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
			status:           models.COMPLETE,
			expectedHTTPCode: 400,
			isValidUUID:      UUIDValidFunc,
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
			status:           models.COMPLETE,
			expectedHTTPCode: 400,
			isValidUUID:      UUIDValidFunc,
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
			status:           models.COMPLETE,
			expectedHTTPCode: 400,
			isValidUUID:      UUIDValidFunc,
		},
		{
			name:             "GET fails with authentification error",
			authIsGiven:      false,
			databaseHasError: false,
			videoAttribute:   "title",
			ascending:        "true",
			page:             "1",
			limit:            "10",
			status:           models.COMPLETE,
			expectedHTTPCode: 401,
			isValidUUID:      UUIDValidFunc,
		},
		{
			name:             "GET fails with database error",
			authIsGiven:      true,
			databaseHasError: true,
			videoAttribute:   "title",
			ascending:        "true",
			page:             "1",
			limit:            "10",
			status:           models.COMPLETE,
			expectedHTTPCode: 500,
			isValidUUID:      UUIDValidFunc,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			// Mock database
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			routerClients := router.Clients{
				UUIDGen: clients.NewUuidGeneratorDummy(nil, tt.isValidUUID),
			}

			// Init videoDAO
			dao_test.ExpectVideosDAOCreation(mock)

			//Create request
			givenRequest := fmt.Sprintf("/api/v1/videos/list/%v/%v/%v/%v/%v", tt.videoAttribute, tt.ascending, tt.page, tt.limit, tt.status.String())

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
				getVideoListQuery := regexp.QuoteMeta(fmt.Sprintf("SELECT * FROM videos WHERE video_status = ? ORDER BY %v %v LIMIT ?,?", tt.videoAttribute, direction))
				getVideoTotal := regexp.QuoteMeta(dao.VideosRequests[dao.GetTotalVideos])

				// Tables
				videosColumns := []string{"id", "title", "video_status", "uploaded_at", "created_at", "updated_at", "source_path", "cover_path"}
				videosRows := sqlmock.NewRows(videosColumns)

				if tt.databaseHasError {
					mock.ExpectQuery(getVideoListQuery).WithArgs(int(tt.status), (pagenum-1)*limitnum, limitnum).WillReturnError(fmt.Errorf("Server Error"))
				} else {
					sourcePathVideo := validVideoId + "/" + "source.mp4"
					coverPath := validVideoId + "/" + "cover.png"
					videosRows.AddRow(validVideoId, "title", int(models.ENCODING), t1, t1, nil, sourcePathVideo, coverPath)
					mock.ExpectQuery(getVideoListQuery).WithArgs(int(tt.status), (pagenum-1)*limitnum, limitnum).WillReturnRows(videosRows)
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
			}, &routerClients, &routerDAO)

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
