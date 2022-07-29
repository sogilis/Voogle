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

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	"github.com/Sogilis/Voogle/src/cmd/api/db/dao_test"
	"github.com/Sogilis/Voogle/src/cmd/api/models"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
	"github.com/Sogilis/Voogle/src/pkg/clients"
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
	sourcePath := validVideoID + "/" + "source.mp4"
	coverPath := validVideoID + "/" + "cover.png"

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
			require.NoError(t, err)
			defer db.Close()

			routerClients := router.Clients{
				UUIDGen: clients.NewUuidGeneratorDummy(nil, tt.isValidUUID),
			}

			dao_test.ExpectVideosDAOCreation(mock)

			if !tt.giveWithAuth || tt.giveRequest == "/api/v1/videos/"+invalidVideoID+"/status" {
				// All these cases will stop before modifying the database : Nothing to do

			} else {
				// Queries
				getVideoFromIdQuery := regexp.QuoteMeta(dao.VideosRequests[dao.GetVideo])

				// Tables
				videosColumns := []string{"id", "title", "video_status", "uploaded_at", "created_at", "updated_at", "source_path", "cover_path"}
				videosRows := sqlmock.NewRows(videosColumns)

				if tt.giveDatabaseErr {
					mock.ExpectQuery(getVideoFromIdQuery).WillReturnError(fmt.Errorf("unknow invalid video ID"))

				} else if tt.giveRequest == "/api/v1/videos/"+unknownVideoID+"/status" {
					mock.ExpectQuery(getVideoFromIdQuery).WillReturnRows(videosRows)

				} else {
					videosRows.AddRow(validVideoID, videoTitle, models.ENCODING, nil, t1, nil, sourcePath, coverPath)
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
			}, &routerClients, &routerDAO)

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
