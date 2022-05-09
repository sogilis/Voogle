package controllers_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
	. "github.com/Sogilis/Voogle/src/cmd/api/router"
)

func TestVideoDelete(t *testing.T) { //nolint:cyclop
	givenUsername := "dev"
	givenUserPwd := "test"

	validVideoID := "1508e7d5-5bc6-4a50-9176-ab0371aa65fe"
	invalidVideoID := "invalidvideoid"
	unknownVideoID := "0000a0a0-0aa0-0a00-0000-aa0000aa00aa"
	UUIDValidFunc := func(u string) bool { _, err := uuid.Parse(u); return err == nil }

	cases := []struct {
		name             string
		giveRequest      string
		giveWithAuth     bool
		giveDatabaseErr  bool
		expectedHTTPCode int
		isValidUUID      func(string) bool
	}{
		{
			name:             "DELETE video ",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/delete",
			giveWithAuth:     true,
			expectedHTTPCode: 200,
			isValidUUID:      UUIDValidFunc},
		{
			name:             "DELETE fails with invalid video ID",
			giveRequest:      "/api/v1/videos/" + invalidVideoID + "/delete",
			giveWithAuth:     true,
			expectedHTTPCode: 400,
			isValidUUID:      UUIDValidFunc},
		{
			name:             "DELETE fails with unknown video ID",
			giveRequest:      "/api/v1/videos/" + unknownVideoID + "/delete",
			giveWithAuth:     true,
			expectedHTTPCode: 404,
			isValidUUID:      UUIDValidFunc},
		{
			name:             "DELETE fails with database error",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/delete",
			giveWithAuth:     true,
			giveDatabaseErr:  true,
			expectedHTTPCode: 500,
			isValidUUID:      UUIDValidFunc},

		{
			name:             "DELETE fails with no auth",
			giveRequest:      "/api/v1/videos/" + validVideoID + "/delete",
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

			routerUUIDGen := UUIDGenerator{
				UUIDGen: uuidgenerator.NewUuidGeneratorDummy(nil, tt.isValidUUID),
			}

			if !tt.giveWithAuth || tt.giveRequest == "/api/v1/videos/"+invalidVideoID+"/delete" {
				// All these cases will stop before modifying the database : Nothing to do

			} else {
				// Queries
				deleteVideo := regexp.QuoteMeta("DELETE FROM videos WHERE id = ?")
				deleteUpload := regexp.QuoteMeta("DELETE FROM uploads WHERE video_id = ?")

				if tt.giveDatabaseErr {
					mock.ExpectPrepare(deleteVideo)
					mock.ExpectExec(deleteVideo).WithArgs(validVideoID).WillReturnError(fmt.Errorf("database internal error"))

				} else if tt.giveRequest == "/api/v1/videos/"+unknownVideoID+"/delete" {
					mock.ExpectPrepare(deleteVideo)
					mock.ExpectExec(deleteVideo).WithArgs(unknownVideoID).WillReturnError(sql.ErrNoRows)

				} else {
					mock.ExpectPrepare(deleteVideo)
					mock.ExpectExec(deleteVideo).WithArgs(validVideoID).WillReturnResult(sqlmock.NewResult(0, 1))

					mock.ExpectPrepare(deleteUpload)
					mock.ExpectExec(deleteUpload).WithArgs(validVideoID).WillReturnResult(sqlmock.NewResult(0, 1))
				}
			}

			r := NewRouter(config.Config{
				UserAuth: givenUsername,
				PwdAuth:  givenUserPwd,
			}, &routerClients, &routerUUIDGen)

			w := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodDelete, tt.giveRequest, nil)
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
