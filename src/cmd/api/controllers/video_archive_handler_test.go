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

func TestVideoArchive(t *testing.T) { //nolint:cyclop
	givenUsername := "dev"
	givenUserPwd := "test"

	validVideoID := "1508e7d5-5bc6-4a50-9176-ab0371aa65fe"
	invalidVideoID := "invalidvideoid"
	unknownVideoID := "0000a0a0-0aa0-0a00-0000-aa0000aa00aa"
	UUIDValidFunc := func(u string) bool { _, err := uuid.Parse(u); return err == nil }
	videoTitle := "title"
	t1 := time.Now()
	sourcePath := validVideoID + "/" + "source.mp4"
	coverPath := validVideoID + "/" + "cover.jpg"

	cases := []struct {
		name             string
		giveRequest      string
		giveWithAuth     bool
		giveDbGetErr     bool
		giveDbUpdateErr  bool
		status           models.VideoStatus
		expectedHTTPCode int
		isValidUUID      func(string) bool
	}{
		{
			name:             "PUT archive video",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/archive",
			giveWithAuth:     true,
			status:           models.COMPLETE,
			expectedHTTPCode: 200,
			isValidUUID:      UUIDValidFunc,
		},
		{
			name:             "PUT fails with status not COMPLETE",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/archive",
			giveWithAuth:     true,
			status:           models.ENCODING,
			expectedHTTPCode: 400,
			isValidUUID:      UUIDValidFunc,
		},
		{
			name:             "PUT fails with database error",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/archive",
			giveWithAuth:     true,
			giveDbGetErr:     true,
			status:           models.COMPLETE,
			expectedHTTPCode: 500,
			isValidUUID:      UUIDValidFunc,
		},
		{
			name:             "PUT fails with unknown video ID",
			giveRequest:      "/api/v1/videos/" + unknownVideoID + "/archive",
			giveWithAuth:     true,
			status:           models.COMPLETE,
			expectedHTTPCode: 404,
			isValidUUID:      UUIDValidFunc,
		},
		{
			name:             "PUT fails with no auth",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/archive",
			giveWithAuth:     false,
			status:           models.COMPLETE,
			expectedHTTPCode: 401,
			isValidUUID:      UUIDValidFunc,
		},
		{
			name:             "PUT fails with invalid video ID",
			giveRequest:      "/api/v1/videos/" + invalidVideoID + "/archive",
			giveWithAuth:     true,
			status:           models.COMPLETE,
			expectedHTTPCode: 400,
			isValidUUID:      UUIDValidFunc,
		},
		{
			name:             "PUT fails with database update fails",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/archive",
			giveWithAuth:     true,
			status:           models.COMPLETE,
			giveDbUpdateErr:  true,
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

			dao_test.ExpectVideosDAOCreation(mock)

			if !tt.giveWithAuth || tt.giveRequest == "/api/v1/videos/"+invalidVideoID+"/archive" {
				// All these cases will stop before modifying the database : Nothing to do

			} else {
				// Queries
				getVideoFromIdQuery := regexp.QuoteMeta(dao.VideosRequests[dao.GetVideo])
				updateVideoQuery := regexp.QuoteMeta(dao.VideosRequests[dao.UpdateVideo])

				// Tables
				videosColumns := []string{"id", "title", "video_status", "uploaded_at", "created_at", "updated_at", "source_path", "cover_path"}
				videosRows := sqlmock.NewRows(videosColumns)

				// Define database response according to case
				if tt.giveDbGetErr {
					mock.ExpectQuery(getVideoFromIdQuery).WillReturnError(fmt.Errorf("unknow invalid video ID"))

				} else if tt.giveRequest == "/api/v1/videos/"+unknownVideoID+"/archive" {
					mock.ExpectQuery(getVideoFromIdQuery).WillReturnRows(videosRows)
				} else {
					videosRows.AddRow(validVideoID, videoTitle, int(tt.status), t1, t1, nil, sourcePath, coverPath)
					mock.ExpectQuery(getVideoFromIdQuery).WillReturnRows(videosRows)

					if tt.status == models.COMPLETE {
						if tt.giveDbUpdateErr {
							mock.ExpectExec(updateVideoQuery).
								WithArgs(videoTitle, int(models.ARCHIVE), t1, sourcePath, coverPath, validVideoID).
								WillReturnError(fmt.Errorf("Error will update db"))
						} else {
							mock.ExpectExec(updateVideoQuery).
								WithArgs(videoTitle, int(models.ARCHIVE), t1, sourcePath, coverPath, validVideoID).
								WillReturnResult(sqlmock.NewResult(0, 1))
						}
					}
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

			req := httptest.NewRequest("PUT", tt.giveRequest, nil)
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
