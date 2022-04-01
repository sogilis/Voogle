package controllers_test

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	"github.com/Sogilis/Voogle/src/pkg/uuidgenerator"

	"github.com/Sogilis/Voogle/src/cmd/api/config"
	"github.com/Sogilis/Voogle/src/cmd/api/models"
	. "github.com/Sogilis/Voogle/src/cmd/api/router"
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

func TestVideoUploadHandler(t *testing.T) {
	givenUsername := "dev"
	givenUserPwd := "test"

	cases := []struct {
		name             string
		giveRequest      string
		giveWithAuth     bool
		giveTitle        string
		giveFieldPart    string
		giveEmptyBody    bool
		giveWrongMagic   bool
		expectedHTTPCode int
		genUUID          func() (string, error)
		putObject        func(io.Reader, string) error
	}{
		{
			name:             "POST upload video",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     true,
			giveTitle:        "title-of-video",
			giveFieldPart:    "video",
			giveWrongMagic:   false,
			expectedHTTPCode: 200,
			genUUID:          func() (string, error) { return "AnUniqueId", nil },
			putObject:        func(f io.Reader, s string) error { _, err := io.ReadAll(f); return err }},
		{
			name:             "POST upload video with space in title",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     true,
			giveTitle:        "title of video",
			giveFieldPart:    "video",
			giveWrongMagic:   false,
			expectedHTTPCode: 200,
			genUUID:          func() (string, error) { return "AnUniqueId", nil },
			putObject: func(f io.Reader, t string) error {
				fmt.Println("title:", t)
				_, err := io.ReadAll(f)
				if strings.Contains(t, " ") {
					return fmt.Errorf("Contains space")
				}
				return err
			}},
		{
			name:             "POST fails with empty title",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     true,
			giveTitle:        "",
			giveFieldPart:    "video",
			expectedHTTPCode: 400,
			genUUID:          func() (string, error) { return "AnUniqueId", nil },
			putObject:        func(f io.Reader, s string) error { _, err := io.ReadAll(f); return err }},
		{
			name:             "POST fails with empty body",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     true,
			giveTitle:        "title-of-video",
			giveFieldPart:    "video",
			giveEmptyBody:    true,
			expectedHTTPCode: 400,
			genUUID:          func() (string, error) { return "AnUniqueId", nil },
			putObject:        func(f io.Reader, s string) error { _, err := io.ReadAll(f); return err }},
		{
			name:             "POST fails with wrong part title",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     true,
			giveTitle:        "title-of-video",
			giveFieldPart:    "vdeo",
			expectedHTTPCode: 400,
			genUUID:          func() (string, error) { return "AnUniqueId", nil },
			putObject:        func(f io.Reader, s string) error { _, err := io.ReadAll(f); return err }},
		{
			name:             "POST fails with wrong magic number",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     true,
			giveTitle:        "title-of-video",
			giveFieldPart:    "video",
			giveWrongMagic:   true,
			expectedHTTPCode: 400,
			genUUID:          func() (string, error) { return "AnUniqueId", nil },
			putObject:        func(f io.Reader, s string) error { _, err := io.ReadAll(f); return err }},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			s3Client := clients.NewS3ClientDummy(nil, nil, tt.putObject, nil)
			amqpClient := clients.NewAmqpClientDummy(nil, nil)

			// Mock database
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			routerClients := Clients{
				S3Client:      s3Client,
				AmqpClient:    amqpClient,
				MariadbClient: db,
			}

			uuidGen := uuidgenerator.NewUuidGeneratorDummy(tt.genUUID)

			routerUUIDGen := UUIDGenerator{
				UUIDGen: uuidGen,
			}

			VideoID, _ := tt.genUUID()
			UploadID, _ := tt.genUUID()

			t1 := time.Now()

			// Create Video
			mock.ExpectExec("INSERT INTO videos").WithArgs(VideoID, tt.giveTitle, int(models.UPLOADING)).WillReturnResult(sqlmock.NewResult(1, 1))

			query := regexp.QuoteMeta("SELECT * FROM videos v WHERE v.id = ?")
			row := sqlmock.NewRows([]string{"id", "title", "video_status", "uploaded_at", "created_at", "updated_at"}).
				AddRow(VideoID, tt.giveTitle, int(models.UPLOADING), nil, t1, t1)
			mock.ExpectQuery(query).WithArgs(VideoID).WillReturnRows(row)

			// Create Upload
			mock.ExpectExec("INSERT INTO uploads").WithArgs(UploadID, VideoID, int(models.STARTED)).WillReturnResult(sqlmock.NewResult(1, 1))

			query = regexp.QuoteMeta("SELECT * FROM uploads u WHERE u.id = ?")
			row = sqlmock.NewRows([]string{"id", "video_id", "upload_status", "uploaded_at", "created_at", "updated_at"}).
				AddRow(UploadID, VideoID, int(models.STARTED), nil, t1, t1)
			mock.ExpectQuery(query).WithArgs(VideoID).WillReturnRows(row)

			// Update videos status : UPLOADED + Upload date
			query = regexp.QuoteMeta("UPDATE videos SET title = ?, video_status = ?, uploaded_at = ?, updated_at = ? WHERE id = ?")
			mock.ExpectExec(query).WithArgs(tt.giveTitle, int(models.UPLOADED), AnyTime{}, t1, VideoID).WillReturnResult(sqlmock.NewResult(0, 1))

			// Update uploads status : DONE + Upload date
			query = regexp.QuoteMeta("UPDATE uploads SET video_id = ?, upload_status = ?, uploaded_at = ?, updated_at = ? WHERE id = ?")
			mock.ExpectExec(query).WithArgs(VideoID, int(models.DONE), AnyTime{}, t1, UploadID).WillReturnResult(sqlmock.NewResult(0, 1))

			// Update video status : ENCODING
			query = regexp.QuoteMeta("UPDATE videos SET title = ?, video_status = ?, uploaded_at = ?, updated_at = ? WHERE id = ?")
			mock.ExpectExec(query).WithArgs(tt.giveTitle, int(models.ENCODING), AnyTime{}, t1, VideoID).WillReturnResult(sqlmock.NewResult(0, 1))

			// Dummy multipart file creation
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			err = writer.WriteField("title", tt.giveTitle)
			assert.NoError(t, err)

			if !tt.giveEmptyBody {
				fileWriter, _ := writer.CreateFormFile(tt.giveFieldPart, "4K.mp4")
				contentFile := bytes.NewBuffer(make([]byte, 0, 1000))
				_, err := io.Copy(fileWriter, contentFile)
				assert.NoError(t, err)

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
					assert.NoError(t, err)
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
					assert.NoError(t, err)
				}
				assert.NoError(t, err)
			}
			writer.Close()

			r := NewRouter(config.Config{
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
			assert.Equal(t, tt.expectedHTTPCode, w.Code)

		})
	}
}
