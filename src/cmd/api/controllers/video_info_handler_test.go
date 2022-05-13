package controllers_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	"github.com/Sogilis/Voogle/src/cmd/api/db/dao_test"
	"github.com/Sogilis/Voogle/src/cmd/api/models"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
)

func TestVideoInfo(t *testing.T) { //nolint:cyclop
	givenUsername := "dev"
	givenUserPwd := "test"

	validVideoID := "1508e7d5-5bc6-4a50-9176-ab0371aa65fe"
	invalidVideoID := "invalidvideoid"
	unknownVideoID := "0000a0a0-0aa0-0a00-0000-aa0000aa00aa"
	UUIDValidFunc := func(u string) bool { _, err := uuid.Parse(u); return err == nil }
	videoTitle := "title"
	t1 := time.Now()
	sourcePath := validVideoID + "/" + "source.mp4"

	cases := []struct {
		name             string
		giveRequest      string
		giveWithAuth     bool
		giveDatabaseErr  bool
		expectedHTTPCode int
		isValidUUID      func(string) bool
	}{
		{
			name:             "GET fails with database error",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/info",
			giveWithAuth:     true,
			giveDatabaseErr:  true,
			expectedHTTPCode: 500,
			isValidUUID:      UUIDValidFunc},

		{
			name:             "GET fails with unknown video ID",
			giveRequest:      "/api/v1/videos/" + unknownVideoID + "/info",
			giveWithAuth:     true,
			expectedHTTPCode: 404,
			isValidUUID:      UUIDValidFunc},

		{
			name:             "GET fails with no auth",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/info",
			giveWithAuth:     false,
			expectedHTTPCode: 401,
			isValidUUID:      UUIDValidFunc},

		{
			name:             "GET fails with invalid video ID",
			giveRequest:      "/api/v1/videos/" + invalidVideoID + "/info",
			giveWithAuth:     true,
			expectedHTTPCode: 400,
			isValidUUID:      UUIDValidFunc},

		{
			name:             "GET video informations",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/info",
			giveWithAuth:     true,
			expectedHTTPCode: 200,
			isValidUUID:      UUIDValidFunc},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			// Mock database
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			routerUUIDGen := router.UUIDGenerator{
				UUIDGen: uuidgenerator.NewUuidGeneratorDummy(nil, tt.isValidUUID),
			}

			dao_test.ExpectVideosDAOCreation(mock)

			if !tt.giveWithAuth || tt.giveRequest == "/api/v1/videos/"+invalidVideoID+"/info" {
				// All these cases will stop before modifying the database : Nothing to do

			} else {
				// Queries
				getVideoFromIdQuery := regexp.QuoteMeta(dao.VideosRequests[dao.GetVideo])

				// Tables
				videosColumns := []string{"id", "title", "video_status", "uploaded_at", "created_at", "updated_at", "source_path"}
				videosRows := sqlmock.NewRows(videosColumns)

				// Define database response according to case
				if tt.giveDatabaseErr {
					mock.ExpectQuery(getVideoFromIdQuery).WillReturnError(fmt.Errorf("unknow invalid video ID"))

				} else if tt.giveRequest == "/api/v1/videos/"+unknownVideoID+"/info" {
					mock.ExpectQuery(getVideoFromIdQuery).WillReturnRows(videosRows)

				} else {
					videosRows.AddRow(validVideoID, videoTitle, int(models.ENCODING), t1, t1, nil, sourcePath)
					mock.ExpectQuery(getVideoFromIdQuery).WillReturnRows(videosRows)
				}
			}

			videoDAO, err := dao.CreateVideosDAO(context.Background(), db)
			require.NoError(t, err)
			routerDAO := router.DAOs{
				VideosDAO: *videoDAO,
			}

			r := router.NewRouter(config.Config{
				UserAuth: givenUsername,
				PwdAuth:  givenUserPwd,
			}, &router.Clients{}, &routerUUIDGen, &routerDAO)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", tt.giveRequest, nil)
			if tt.giveWithAuth {
				req.SetBasicAuth(givenUsername, givenUserPwd)
			}

			r.ServeHTTP(w, req)
			require.Equal(t, tt.expectedHTTPCode, w.Code)

			// we make sure that all expectations were met
			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})

	}

}
