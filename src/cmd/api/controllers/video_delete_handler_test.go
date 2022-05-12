package controllers_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
)

func TestVideoDelete(t *testing.T) { //nolint:cyclop
	givenUsername := "dev"
	givenUserPwd := "test"

	validVideoID := "1508e7d5-5bc6-4a50-9176-ab0371aa65fe"
	invalidVideoID := "invalidvideoid"
	unknownVideoID := "0000a0a0-0aa0-0a00-0000-aa0000aa00aa"
	UUIDValidFunc := func(u string) bool { _, err := uuid.Parse(u); return err == nil }
	removeObjectS3 := func(v string) error { return nil }

	videoTitle := "title"
	t1 := time.Now()
	sourcePath := validVideoID + "/" + "source.mp4"

	cases := []struct {
		name                 string
		giveRequest          string
		giveWithAuth         bool
		giveDatabaseErr      bool
		expectedHTTPCode     int
		videoDeletionFails   bool
		uploadsDeletionFails bool
		isValidUUID          func(string) bool
		removeObject         func(string) error
	}{
		{
			name:             "DELETE video ",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/delete",
			giveWithAuth:     true,
			expectedHTTPCode: 200,
			isValidUUID:      UUIDValidFunc,
			removeObject:     removeObjectS3},
		{
			name:             "DELETE fails with invalid video ID",
			giveRequest:      "/api/v1/videos/" + invalidVideoID + "/delete",
			giveWithAuth:     true,
			expectedHTTPCode: 400,
			isValidUUID:      UUIDValidFunc,
			removeObject:     removeObjectS3},
		{
			name:             "DELETE fails with unknown video ID",
			giveRequest:      "/api/v1/videos/" + unknownVideoID + "/delete",
			giveWithAuth:     true,
			expectedHTTPCode: 404,
			isValidUUID:      UUIDValidFunc,
			removeObject:     removeObjectS3},
		{
			name:             "DELETE fails with database error",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/delete",
			giveWithAuth:     true,
			giveDatabaseErr:  true,
			expectedHTTPCode: 500,
			isValidUUID:      UUIDValidFunc,
			removeObject:     removeObjectS3},
		{
			name:             "DELETE fails with no auth",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/delete",
			giveWithAuth:     false,
			expectedHTTPCode: 401,
			isValidUUID:      UUIDValidFunc,
			removeObject:     removeObjectS3},
		{
			name:             "DELETE fails with S3 error",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/delete",
			giveWithAuth:     true,
			expectedHTTPCode: 500,
			isValidUUID:      UUIDValidFunc,
			removeObject:     func(v string) error { return fmt.Errorf("S3 error") },
		},
		{
			name:               "DELETE fails with video deletion fails",
			giveRequest:        "/api/v1/videos/" + validVideoID + "/delete",
			giveWithAuth:       true,
			expectedHTTPCode:   500,
			isValidUUID:        UUIDValidFunc,
			removeObject:       removeObjectS3,
			videoDeletionFails: true,
		},
		{
			name:                 "DELETE fails with uploads deletion fails",
			giveRequest:          "/api/v1/videos/" + validVideoID + "/delete",
			giveWithAuth:         true,
			expectedHTTPCode:     500,
			isValidUUID:          UUIDValidFunc,
			removeObject:         removeObjectS3,
			uploadsDeletionFails: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			s3Client := clients.NewS3ClientDummy(nil, nil, nil, nil, tt.removeObject)

			// Mock database
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			routerClients := router.Clients{
				S3Client: s3Client,
			}

			routerUUIDGen := router.UUIDGenerator{
				UUIDGen: uuidgenerator.NewUuidGeneratorDummy(nil, tt.isValidUUID),
			}

			dao.ExpectVideosDAOCreation(mock)
			dao.ExpectUploadsDAOCreation(mock)

			if !tt.giveWithAuth || tt.giveRequest == "/api/v1/videos/"+invalidVideoID+"/delete" {
				// All these cases will stop before modifying the database : Nothing to do

			} else {
				// Queries
				deleteVideo := regexp.QuoteMeta(dao.VideosRequests[dao.DeleteVideo])

				deleteUpload := regexp.QuoteMeta(dao.UploadsRequests[dao.DeleteUpload])
				getVideoFromIdQuery := regexp.QuoteMeta(dao.VideosRequests[dao.GetVideo])

				// Tables
				videosColumns := []string{"id", "title", "video_status", "uploaded_at", "created_at", "updated_at", "source_path"}
				videosRows := sqlmock.NewRows(videosColumns)

				if tt.giveDatabaseErr {
					mock.ExpectQuery(getVideoFromIdQuery).WillReturnError(fmt.Errorf("database internal error"))

				} else if tt.giveRequest == "/api/v1/videos/"+unknownVideoID+"/delete" {
					mock.ExpectQuery(getVideoFromIdQuery).WillReturnRows(videosRows)

				} else {
					videosRows.AddRow(validVideoID, videoTitle, contracts.Video_VIDEO_STATUS_ENCODING, t1, t1, nil, sourcePath)
					mock.ExpectQuery(getVideoFromIdQuery).WillReturnRows(videosRows)

					mock.ExpectBegin()
					if tt.uploadsDeletionFails {
						mock.ExpectExec(deleteUpload).WithArgs(validVideoID).WillReturnError(fmt.Errorf("database internal error"))
						mock.ExpectRollback()

					} else {
						mock.ExpectExec(deleteUpload).WithArgs(validVideoID).WillReturnResult(sqlmock.NewResult(0, 1))

						if tt.videoDeletionFails {
							mock.ExpectExec(deleteVideo).WithArgs(validVideoID).WillReturnError(fmt.Errorf("database internal error"))
							mock.ExpectRollback()

						} else {
							mock.ExpectExec(deleteVideo).WithArgs(validVideoID).WillReturnResult(sqlmock.NewResult(0, 1))
							mock.ExpectCommit()
						}
					}
				}
			}

			videosDAO, err := dao.CreateVideosDAO(context.Background(), db)
			require.NoError(t, err)

			uploadsDAO, err := dao.CreateUploadsDAO(context.Background(), db)
			require.NoError(t, err)

			routerDAO := router.DAOs{
				VideosDAO:  *videosDAO,
				UploadsDAO: *uploadsDAO,
			}

			r := router.NewRouter(config.Config{
				UserAuth: givenUsername,
				PwdAuth:  givenUserPwd,
			}, &routerClients, &routerUUIDGen, &routerDAO)

			w := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodDelete, tt.giveRequest, nil)
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
