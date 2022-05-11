package controllers_test

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/models"
	"github.com/Sogilis/Voogle/src/cmd/api/router"
)

// Used to mock upload_at time.Time that is set into
// video_upload_handler.go.
// See https://github.com/DATA-DOG/go-sqlmock#matching-arguments-like-timetime
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestVideoUploadHandler(t *testing.T) { //nolint:cyclop
	givenUsername := "dev"
	givenUserPwd := "test"

	cases := []struct {
		name               string
		giveRequest        string
		giveWithAuth       bool
		giveTitle          string
		giveFieldPart      string
		giveEmptyBody      bool
		giveWrongMagic     bool
		lastUploadFailed   bool
		lastEncodeFailed   bool
		titleAlreadyExists bool
		createVideoFail    bool
		createUploadFail   bool
		uploadOnS3fail     bool
		expectedHTTPCode   int
		genUUID            func() (string, error)
		putObject          func(io.Reader, string) error
	}{
		{
			name:             "POST upload video",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     true,
			giveTitle:        "title-of-video",
			giveFieldPart:    "video",
			expectedHTTPCode: 200,
			genUUID:          func() (string, error) { return "AUniqueId", nil },
			putObject:        func(f io.Reader, s string) error { _, err := io.ReadAll(f); return err }},
		{
			name:             "POST upload with last video upload failed",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     true,
			giveTitle:        "title-of-video",
			giveFieldPart:    "video",
			lastUploadFailed: true,
			expectedHTTPCode: 200,
			genUUID:          func() (string, error) { return "AUniqueId", nil },
			putObject:        func(f io.Reader, s string) error { _, err := io.ReadAll(f); return err }},
		{
			name:             "POST upload with last video encode failed",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     true,
			giveTitle:        "title-of-video",
			giveFieldPart:    "video",
			lastEncodeFailed: true,
			expectedHTTPCode: 200,
			genUUID:          func() (string, error) { return "AUniqueId", nil },
			putObject:        func(f io.Reader, s string) error { _, err := io.ReadAll(f); return err }},
		{
			name:             "POST fails with empty title",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     true,
			giveTitle:        "",
			giveFieldPart:    "video",
			expectedHTTPCode: 400,
			genUUID:          func() (string, error) { return "AUniqueId", nil },
			putObject:        func(f io.Reader, s string) error { _, err := io.ReadAll(f); return err }},
		{
			name:             "POST fails with empty body",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     true,
			giveTitle:        "title-of-video",
			giveFieldPart:    "video",
			giveEmptyBody:    true,
			expectedHTTPCode: 400,
			genUUID:          func() (string, error) { return "AUniqueId", nil },
			putObject:        func(f io.Reader, s string) error { _, err := io.ReadAll(f); return err }},
		{
			name:             "POST fails with wrong part field",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     true,
			giveTitle:        "title-of-video",
			giveFieldPart:    "NOT-video",
			expectedHTTPCode: 400,
			genUUID:          func() (string, error) { return "AUniqueId", nil },
			putObject:        func(f io.Reader, s string) error { _, err := io.ReadAll(f); return err }},
		{
			name:             "POST fails with wrong magic number",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     true,
			giveTitle:        "title-of-video",
			giveFieldPart:    "video",
			giveWrongMagic:   true,
			expectedHTTPCode: 415,
			genUUID:          func() (string, error) { return "AUniqueId", nil },
			putObject:        func(f io.Reader, s string) error { _, err := io.ReadAll(f); return err }},
		{
			name:               "POST fails with title already exist",
			giveRequest:        "/api/v1/videos/upload",
			giveWithAuth:       true,
			giveTitle:          "title-of-video",
			giveFieldPart:      "video",
			titleAlreadyExists: true,
			expectedHTTPCode:   409,
			genUUID:            func() (string, error) { return "AUniqueId", nil },
			putObject:          func(f io.Reader, s string) error { _, err := io.ReadAll(f); return err }},
		{
			name:             "POST fails with create video fail",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     true,
			giveTitle:        "title-of-video",
			giveFieldPart:    "video",
			createVideoFail:  true,
			expectedHTTPCode: 500,
			genUUID:          func() (string, error) { return "AUniqueId", nil },
			putObject:        func(f io.Reader, s string) error { _, err := io.ReadAll(f); return err }},
		{
			name:             "POST fails with create upload fail",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     true,
			giveTitle:        "title-of-video",
			giveFieldPart:    "video",
			createUploadFail: true,
			expectedHTTPCode: 500,
			genUUID:          func() (string, error) { return "AUniqueId", nil },
			putObject:        func(f io.Reader, s string) error { _, err := io.ReadAll(f); return err }},
		{
			name:             "POST fails with S3 upload failed",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     true,
			giveTitle:        "title-of-video",
			giveFieldPart:    "video",
			expectedHTTPCode: 500,
			uploadOnS3fail:   true,
			genUUID:          func() (string, error) { return "AUniqueId", nil },
			putObject:        func(f io.Reader, s string) error { return errors.New("Cannot upload on S3") }},
		{
			name:             "POST fails with no auth",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     false,
			expectedHTTPCode: 401},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			s3Client := clients.NewS3ClientDummy(nil, nil, tt.putObject, nil, nil)
			amqpClient := clients.NewAmqpClientDummy(nil, nil, nil)

			// Mock database
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			routerClients := router.Clients{
				S3Client:      s3Client,
				AmqpClient:    amqpClient,
				MariadbClient: db,
			}

			routerUUIDGen := router.UUIDGenerator{
				UUIDGen: uuidgenerator.NewUuidGeneratorDummy(tt.genUUID, nil),
			}

			if tt.giveTitle == "" || tt.giveEmptyBody || tt.giveFieldPart == "NOT-video" || tt.giveWrongMagic || !tt.giveWithAuth {
				// All these cases will stop before modifying the database : Nothing to do

			} else {
				// Queries
				createVideoQuery := regexp.QuoteMeta("INSERT INTO videos")
				updateVideoQuery := regexp.QuoteMeta("UPDATE videos SET title = ?, video_status = ?, uploaded_at = ? WHERE id = ?")
				getVideoFromTitleQuery := regexp.QuoteMeta("SELECT * FROM videos v WHERE v.title = ?")
				getVideoFromIdQuery := regexp.QuoteMeta("SELECT * FROM videos v WHERE v.id = ?")

				createUploadQuery := regexp.QuoteMeta("INSERT INTO uploads")
				updateUploadQuery := regexp.QuoteMeta("UPDATE uploads SET video_id = ?, upload_status = ?, uploaded_at = ? WHERE id = ?")
				getUploadQuery := regexp.QuoteMeta("SELECT * FROM uploads u WHERE u.id = ?")

				// Tables
				videosColumns := []string{"id", "title", "video_status", "uploaded_at", "created_at", "updated_at", "source_path"}
				uploadsColumns := []string{"id", "video_id", "upload_status", "uploaded_at", "created_at", "updated_at"}
				videosRows := sqlmock.NewRows(videosColumns)
				uploadRows := sqlmock.NewRows(uploadsColumns)

				VideoID, _ := tt.genUUID()
				UploadID, _ := tt.genUUID()

				t1 := time.Now()
				sourcePath := VideoID + "/" + "source.mp4"

				if tt.titleAlreadyExists {
					// Create Video (fail)
					mock.ExpectPrepare(createVideoQuery)
					mock.ExpectExec(createVideoQuery).
						WithArgs(VideoID, tt.giveTitle, models.UPLOADING, sourcePath).
						WillReturnError(fmt.Errorf("Error while creating new video"))

					videosRows.AddRow(VideoID, tt.giveTitle, models.UPLOADING, nil, t1, t1, sourcePath)
					mock.ExpectPrepare(getVideoFromTitleQuery)
					mock.ExpectQuery(getVideoFromTitleQuery).WithArgs(tt.giveTitle).WillReturnRows(videosRows)

				} else if tt.createVideoFail {
					// Create Video (fail)
					mock.ExpectPrepare(createVideoQuery)
					mock.ExpectExec(createVideoQuery).
						WithArgs(VideoID, tt.giveTitle, models.UPLOADING, sourcePath).
						WillReturnError(fmt.Errorf("Error while creating new video"))

					mock.ExpectPrepare(getVideoFromTitleQuery)
					mock.ExpectQuery(getVideoFromTitleQuery).WithArgs(tt.giveTitle).WillReturnRows(videosRows)

				} else if tt.lastEncodeFailed {
					// Create Video (fail)
					mock.ExpectPrepare(createVideoQuery)
					mock.ExpectExec(createVideoQuery).
						WithArgs(VideoID, tt.giveTitle, models.UPLOADING, sourcePath).
						WillReturnError(fmt.Errorf("Duplicate entry : 1062"))

					videosRows.AddRow(VideoID, tt.giveTitle, models.FAIL_ENCODE, nil, t1, t1, sourcePath)
					mock.ExpectPrepare(getVideoFromTitleQuery)
					mock.ExpectQuery(getVideoFromTitleQuery).WithArgs(tt.giveTitle).WillReturnRows(videosRows)

					// Update video status : ENCODING
					mock.ExpectPrepare(updateVideoQuery)
					mock.ExpectExec(updateVideoQuery).
						WithArgs(tt.giveTitle, models.ENCODING, nil, VideoID).
						WillReturnResult(sqlmock.NewResult(0, 1))

				} else {
					if tt.lastUploadFailed {
						// Create Video (fail)
						mock.ExpectPrepare(createVideoQuery)
						mock.ExpectExec(createVideoQuery).
							WithArgs(VideoID, tt.giveTitle, models.UPLOADING, sourcePath).
							WillReturnError(fmt.Errorf("Duplicate entry : 1062"))

						videosRows.AddRow(VideoID, tt.giveTitle, models.FAIL_UPLOAD, nil, t1, t1, sourcePath)
						mock.ExpectPrepare(getVideoFromTitleQuery)
						mock.ExpectQuery(getVideoFromTitleQuery).WithArgs(tt.giveTitle).WillReturnRows(videosRows)

					} else {
						// Create Video
						mock.ExpectPrepare(createVideoQuery)
						mock.ExpectExec(createVideoQuery).
							WithArgs(VideoID, tt.giveTitle, models.UPLOADING, sourcePath).
							WillReturnResult(sqlmock.NewResult(1, 1))

						videosRows.AddRow(VideoID, tt.giveTitle, models.UPLOADING, nil, t1, t1, sourcePath)
						mock.ExpectPrepare(getVideoFromIdQuery)
						mock.ExpectQuery(getVideoFromIdQuery).WithArgs(VideoID).WillReturnRows(videosRows)
					}

					if tt.createUploadFail {
						// Create Upload (fail)
						mock.ExpectPrepare(createUploadQuery)
						mock.ExpectExec(createUploadQuery).
							WithArgs(UploadID, VideoID, models.STARTED).
							WillReturnError(fmt.Errorf("Error while creating new upload"))

						// Update videos status : FAIL_UPLOAD
						mock.ExpectPrepare(updateVideoQuery)
						mock.ExpectExec(updateVideoQuery).
							WithArgs(tt.giveTitle, models.FAIL_UPLOAD, nil, VideoID).
							WillReturnResult(sqlmock.NewResult(0, 1))

					} else {
						// Create Upload
						mock.ExpectPrepare(createUploadQuery)
						mock.ExpectExec(createUploadQuery).
							WithArgs(UploadID, VideoID, models.STARTED).
							WillReturnResult(sqlmock.NewResult(1, 1))

						uploadRows.AddRow(UploadID, VideoID, models.STARTED, nil, t1, t1)
						mock.ExpectPrepare(getUploadQuery)
						mock.ExpectQuery(getUploadQuery).WithArgs(VideoID).WillReturnRows(uploadRows)

						if tt.uploadOnS3fail {
							mock.ExpectBegin()

							mock.ExpectPrepare(updateVideoQuery)
							mock.ExpectExec(updateVideoQuery).
								WithArgs(tt.giveTitle, models.FAIL_UPLOAD, nil, VideoID).
								WillReturnResult(sqlmock.NewResult(0, 1))

							mock.ExpectPrepare(updateUploadQuery)
							mock.ExpectExec(updateUploadQuery).
								WithArgs(VideoID, models.FAILED, nil, UploadID).
								WillReturnResult(sqlmock.NewResult(0, 1))

							mock.ExpectCommit()

						} else {
							// Update videos status : UPLOADED + Upload date
							mock.ExpectPrepare(updateVideoQuery)
							mock.ExpectExec(updateVideoQuery).
								WithArgs(tt.giveTitle, models.UPLOADED, AnyTime{}, VideoID).
								WillReturnResult(sqlmock.NewResult(0, 1))

							// Update uploads status : DONE + Upload date
							mock.ExpectPrepare(updateUploadQuery)
							mock.ExpectExec(updateUploadQuery).
								WithArgs(VideoID, models.DONE, AnyTime{}, UploadID).
								WillReturnResult(sqlmock.NewResult(0, 1))

							// Update video status : ENCODING
							mock.ExpectPrepare(updateVideoQuery)
							mock.ExpectExec(updateVideoQuery).
								WithArgs(tt.giveTitle, models.ENCODING, AnyTime{}, VideoID).
								WillReturnResult(sqlmock.NewResult(0, 1))
						}
					}
				}
			}

			// Dummy multipart file creation
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			err = writer.WriteField("title", tt.giveTitle)
			require.NoError(t, err)

			if !tt.giveEmptyBody {
				fileWriter, _ := writer.CreateFormFile(tt.giveFieldPart, "4K.mp4")
				contentFile := bytes.NewBuffer(make([]byte, 0, 1000))
				_, err := io.Copy(fileWriter, contentFile)
				require.NoError(t, err)

				if !tt.giveWrongMagic {
					// Webm magic number
					data := []byte{
						0x1a, 0x45, 0xdf, 0xa3, 0x9f, 0x42, 0x86, 0x81, 0x01, 0x42, 0xf7, 0x81, 0x01, 0x42, 0xf2, 0x81,
						0x04, 0x42, 0xf3, 0x81, 0x08, 0x42, 0x82, 0x84, 0x77, 0x65, 0x62, 0x6d, 0x42, 0x87, 0x81, 0x02,
						0x42, 0x85, 0x81, 0x02, 0x18, 0x53, 0x80, 0x67, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x4a, 0xf7,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					}

					_, err := fileWriter.Write(data)
					require.NoError(t, err)
				} else {
					// Wrong magic number
					data := []byte{
						0x15, 0x00, 0xe4, 0xaf, 0x0f, 0x42, 0x86, 0x81, 0x01, 0x42, 0xf7, 0x81, 0x01, 0x42, 0xf2, 0x81,
						0x04, 0x42, 0xf3, 0x81, 0x08, 0x42, 0x82, 0x84, 0x77, 0x65, 0x62, 0x6d, 0x42, 0x87, 0x81, 0x02,
						0x42, 0x85, 0x81, 0x02, 0x18, 0x53, 0x80, 0x67, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x4a, 0xf7,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					}

					_, err := fileWriter.Write(data)
					require.NoError(t, err)
				}
				require.NoError(t, err)
			}
			writer.Close()

			r := router.NewRouter(config.Config{
				UserAuth: givenUsername,
				PwdAuth:  givenUserPwd,
			}, &routerClients, &routerUUIDGen)

			w := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodPost, tt.giveRequest, body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
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
