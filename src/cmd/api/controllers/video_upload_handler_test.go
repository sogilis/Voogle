package controllers_test

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Sogilis/Voogle/src/cmd/api/clients"
	"github.com/Sogilis/Voogle/src/cmd/api/config"
	. "github.com/Sogilis/Voogle/src/cmd/api/router"
)

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
		expectedHTTPCode int
		putObject        func(io.Reader, string) error
	}{
		{
			name:             "POST upload video",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     true,
			giveTitle:        "title-of-video",
			giveFieldPart:    "video",
			expectedHTTPCode: 200,
			putObject:        func(f io.Reader, s string) error { _, err := io.ReadAll(f); return err }},
		{
			name:             "POST upload video with space in title",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     true,
			giveTitle:        "title of video",
			giveFieldPart:    "video",
			expectedHTTPCode: 200,
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
			putObject:        func(f io.Reader, s string) error { _, err := io.ReadAll(f); return err }},
		{
			name:             "POST fails with empty body",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     true,
			giveTitle:        "title-of-video",
			giveFieldPart:    "video",
			giveEmptyBody:    true,
			expectedHTTPCode: 400,
			putObject:        func(f io.Reader, s string) error { _, err := io.ReadAll(f); return err }},
		{
			name:             "POST fails with wrong part title",
			giveRequest:      "/api/v1/videos/upload",
			giveWithAuth:     true,
			giveTitle:        "title-of-video",
			giveFieldPart:    "vdeo",
			expectedHTTPCode: 400,
			putObject:        func(f io.Reader, s string) error { _, err := io.ReadAll(f); return err }},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			s3Client := clients.NewS3ClientDummy(nil, nil, tt.putObject)
			redisClient := clients.NewRedisClientDummy(nil, nil, nil)

			routerClients := Clients{
				S3Client:    s3Client,
				RedisClient: redisClient,
			}

			// Dummy multipart file creation
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			err := writer.WriteField("title", tt.giveTitle)
			assert.Nil(t, err)

			if !tt.giveEmptyBody {
				fileWriter, _ := writer.CreateFormFile(tt.giveFieldPart, "4K.mp4")
				contentFile := bytes.NewBuffer(make([]byte, 0, 1000))
				_, err := io.Copy(fileWriter, contentFile)
				assert.Nil(t, err)
			}
			writer.Close()

			r := NewRouter(config.Config{
				UserAuth: givenUsername,
				PwdAuth:  givenUserPwd,
			}, &routerClients)

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
