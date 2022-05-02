package end2end_tests

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/franela/goblin"
	"github.com/stretchr/testify/assert"

	"github.com/Sogilis/Voogle/end2end/helpers"
)

func Test_Videos(t *testing.T) {
	host := os.Getenv("E2E_API_ENDPOINT")

	user := os.Getenv("E2E_USER_NAME")
	pwd := os.Getenv("E2E_USER_PWD")

	g := Goblin(t)
	g.Describe("Videos >", func() {
		g.Before(func() {
			// TODO clean DATA
		})

		g.Describe("List >", func() {
			path := "/api/v1/videos/list/title/true/1/10"

			g.Describe("Without login >", func() {
				session := helpers.NewSession(host)
				g.It("Returns a 401", func() {
					t.Log("PATH - GET - " + path)

					code, _, err := session.Get(path)
					assert.NoError(t, err)

					assert.Equal(t, 401, code)
				})
			})

			g.Describe("With login >", func() {
				session := helpers.NewSession(host)
				g.Before(func() {
					assert.Nil(t, session.Login(user, pwd))
				})
				g.It("Returns an empty list of videos", func() {
					t.Log("PATH - GET - " + path)

					code, body, err := session.Get(path)
					assert.NoError(t, err)

					assert.Equal(t, 200, code)

					// Reading the body
					rawBody, err := ioutil.ReadAll(body)
					g.Assert(err).IsNil()
					var videoData helpers.VideoListResponse
					err = json.Unmarshal(rawBody, &videoData)
					assert.NoError(t, err)

					assert.Equal(t, 1, videoData.LastPage)
					assert.Equal(t, 0, len(videoData.Videos))
				})

				g.It("Returns a list of videos with One element", func() {
					g.Timeout(time.Duration(60) * time.Second)
					t.Log("PATH - GET - " + path)

					f, err := os.Open("../samples/1280x720_2mb.mp4")
					assert.NoError(t, err)
					defer f.Close()

					// Post video upload
					code, body, err := session.PostMultipart("/api/v1/videos/upload", "test data", "video.avi", f)
					assert.NoError(t, err)

					assert.Equal(t, 200, code)

					// Reading the body
					rawBody, err := ioutil.ReadAll(body)
					g.Assert(err).IsNil()
					var uploadResponse helpers.Response
					err = json.Unmarshal(rawBody, &uploadResponse)
					assert.NoError(t, err)

					// Get video status
					var videoStatus helpers.VideoStatus
					for strings.ToLower(videoStatus.Status) != "complete" {
						time.Sleep(5 * time.Second)
						code, body, err = session.Get("/api/v1/videos/" + uploadResponse.Video.ID + "/status")
						assert.NoError(t, err)

						assert.Equal(t, 200, code)

						// Reading the body
						rawBody, err = ioutil.ReadAll(body)
						g.Assert(err).IsNil()
						err = json.Unmarshal(rawBody, &videoStatus)
						assert.NoError(t, err)
					}

					// Get video list
					code, body, err = session.Get(path)
					assert.NoError(t, err)

					assert.Equal(t, 200, code)

					// Reading the body
					rawBody, err = ioutil.ReadAll(body)
					g.Assert(err).IsNil()
					var videoData helpers.VideoListResponse
					err = json.Unmarshal(rawBody, &videoData)
					assert.NoError(t, err)

					assert.Equal(t, 1, videoData.LastPage)
					assert.Equal(t, 1, len(videoData.Videos))
					assert.Equal(t, "test data", videoData.Videos[0].Title)
				})

				g.It("Returns an error title already exist", func() {
					t.Log("PATH - GET - " + path)

					f, err := os.Open("../samples/1280x720_2mb.mp4")
					assert.NoError(t, err)
					defer f.Close()

					// Post video upload
					code, _, err := session.PostMultipart("/api/v1/videos/upload", "test data", "video.avi", f)
					assert.NoError(t, err)

					assert.Equal(t, 409, code)
				})

				g.It("Returns video stream master and first part", func() {
					// Get video list
					code, body, err := session.Get(path)
					assert.NoError(t, err)

					assert.Equal(t, 200, code)

					// Reading the body
					rawBody, err := ioutil.ReadAll(body)
					g.Assert(err).IsNil()
					var videoData helpers.AllVideos
					err = json.Unmarshal(rawBody, &videoData)
					assert.NoError(t, err)

					assert.Equal(t, "Success", videoData.Status)
					assert.Equal(t, 1, len(videoData.Data))
					assert.Equal(t, "test data", videoData.Data[0].Title)

					// Get video master
					code, body, err = session.Get("/api/v1/videos/" + videoData.Data[0].Id + "/streams/master.m3u8")
					assert.NoError(t, err)

					assert.Equal(t, 200, code)

					// Reading the body
					rawBody, err = ioutil.ReadAll(body)
					g.Assert(err).IsNil()

					//Parse response
					parseBody := strings.Split(string(rawBody), "#")
					assert.Equal(t, "EXTM3U\n", parseBody[1])
					assert.Equal(t, "EXT-X-VERSION:3\n", parseBody[2])

					// Get video part
					code, _, err = session.Get("/api/v1/videos/" + videoData.Data[0].Id + "/streams/v0/segment0.ts")
					assert.NoError(t, err)

					assert.Equal(t, 200, code)
				})
			})
		})
	})
}
