package integration_tests

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	. "github.com/franela/goblin"
	"github.com/stretchr/testify/require"

	"github.com/Sogilis/Voogle/integration/helpers"
)

func Test_Videos(t *testing.T) {
	host := os.Getenv("INTEGRATION_API_ENDPOINT")
	user := os.Getenv("INTEGRATION_USER_NAME")
	pwd := os.Getenv("INTEGRATION_USER_PWD")

	g := Goblin(t)
	g.Describe("Videos >", func() {
		sessionNil := helpers.NewSession(host)

		session := helpers.NewSession(host)
		require.Nil(t, session.Login(user, pwd))

		pathUpload := "/api/v1/videos/upload"
		pathList := "/api/v1/videos/list/title/true/1/10/Complete"

		videoTitle := "test"
		var videoID string

		g.Describe("Upload >", func() {
			g.Before(func() {
				// Clear data
				_, _ = session.Delete("/api/v1/videos/" + videoID + "/delete")
			})

			g.It("Upload one video", func() {
				t.Log("PATH - POST - " + pathUpload)

				// Open video file
				f, err := os.Open("../samples/1280x720_2mb.mp4")
				require.NoError(t, err)
				defer f.Close()

				// Post video upload
				code, body, err := session.PostMultipart("/api/v1/videos/upload", videoTitle, "video.avi", f)
				require.NoError(t, err)
				require.Equal(t, 200, code)

				// Reading the body
				rawBody, err := ioutil.ReadAll(body)
				require.NoError(t, err)

				// Retrieve video informations
				var uploadResponse helpers.Response
				err = json.Unmarshal(rawBody, &uploadResponse)
				require.NoError(t, err)
				require.Equal(t, uploadResponse.Video.Title, videoTitle)

				// Update videoID
				videoID = uploadResponse.Video.ID
			})

			g.It("Returns an error title already exist", func() {
				t.Log("PATH - POST - " + pathUpload)

				// Open video file
				f, err := os.Open("../samples/1280x720_2mb.mp4")
				require.NoError(t, err)
				defer f.Close()

				// Post video upload with same title
				code, _, err := session.PostMultipart("/api/v1/videos/upload", videoTitle, "video.avi", f)
				require.NoError(t, err)
				require.Equal(t, 409, code)
			})

			g.It("Returns an error unsported media format", func() {
				t.Log("PATH - POST - " + pathUpload)

				// Open image file
				fImage, err := os.Open("../samples/image.mp4")
				require.NoError(t, err)
				defer fImage.Close()

				// Post video upload with image file
				code, _, err := session.PostMultipart("/api/v1/videos/upload", "test image", "falsevideo.mp4", fImage)
				require.NoError(t, err)
				require.Equal(t, 415, code)
			})
		})

		g.Describe("List >", func() {
			g.Describe("Without login >", func() {
				g.It("Returns a 401", func() {
					t.Log("PATH - GET - " + pathList)
					code, _, err := sessionNil.Get(pathList)
					require.NoError(t, err)
					require.Equal(t, 401, code)
				})
			})

			g.Describe("With login >", func() {
				g.Before(func() {
					// Clear data
					_, _ = session.Delete("/api/v1/videos/" + videoID + "/delete")
				})

				g.Describe("With empty list >", func() {
					g.It("Returns an empty list of videos", func() {
						t.Log("PATH - GET - " + pathList)

						// Get video list
						code, body, err := session.Get(pathList)
						require.NoError(t, err)
						require.Equal(t, 200, code)

						// Reading the body
						rawBody, err := ioutil.ReadAll(body)
						require.NoError(t, err)
						var videoData helpers.VideoListResponse
						err = json.Unmarshal(rawBody, &videoData)
						require.NoError(t, err)

						require.Equal(t, 0, len(videoData.Videos))
					})
				})

				g.Describe("With one video >", func() {
					g.Before(func() {
						// Open video file
						f, err := os.Open("../samples/1280x720_2mb.mp4")
						require.NoError(t, err)
						defer f.Close()

						// Post video upload
						_, body, _ := session.PostMultipart(pathUpload, videoTitle, "video.avi", f)

						// Reading the body
						rawBody, _ := ioutil.ReadAll(body)
						var uploadResponse helpers.Response
						_ = json.Unmarshal(rawBody, &uploadResponse)

						// Update videoID and Get video status
						videoID = uploadResponse.Video.ID
						_ = session.WaitVideoEncoded("/api/v1/videos/" + videoID + "/status")
					})

					g.It("Returns a list of videos with One element", func() {
						t.Log("PATH - GET - " + pathList)

						// Get video list
						code, body, err := session.Get(pathList)
						require.NoError(t, err)
						require.Equal(t, 200, code)

						// Reading the body
						rawBody, err := ioutil.ReadAll(body)
						require.NoError(t, err)
						var videoData helpers.VideoListResponse
						err = json.Unmarshal(rawBody, &videoData)
						require.NoError(t, err)

						require.Equal(t, 1, len(videoData.Videos))
						require.Equal(t, videoTitle, videoData.Videos[0].Title)
					})
				})
			})
		})

		g.Describe("Stream >", func() {
			g.Before(func() {
				// Clear data
				_, _ = session.Delete("/api/v1/videos/" + videoID + "/delete")

				// Open video file
				f, err := os.Open("../samples/1280x720_2mb.mp4")
				require.NoError(t, err)
				defer f.Close()

				// Post video upload
				_, body, _ := session.PostMultipart(pathUpload, videoTitle, "video.avi", f)

				// Reading the body
				rawBody, _ := ioutil.ReadAll(body)
				var uploadResponse helpers.Response
				_ = json.Unmarshal(rawBody, &uploadResponse)

				// Update videoID and Get video status
				videoID = uploadResponse.Video.ID
				_ = session.WaitVideoEncoded("/api/v1/videos/" + videoID + "/status")
			})

			g.It("Returns video stream master and first part", func() {
				// Get video master
				code, body, err := session.Get("/api/v1/videos/" + videoID + "/streams/master.m3u8")
				require.NoError(t, err)
				require.Equal(t, 200, code)

				// Reading the body
				rawBody, err := ioutil.ReadAll(body)
				require.NoError(t, err)

				// Check response
				parseBody := strings.Split(string(rawBody), "#")
				require.Equal(t, "EXTM3U\n", parseBody[1])
				require.Equal(t, "EXT-X-VERSION:3\n", parseBody[2])

				// Get video part
				code, _, err = session.Get("/api/v1/videos/" + videoID + "/streams/v0/segment0.ts")
				require.NoError(t, err)
				require.Equal(t, 200, code)
			})
		})

		g.Describe("Delete >", func() {
			g.Before(func() {
				// Clear data
				_, _ = session.Delete("/api/v1/videos/" + videoID + "/delete")

				// Open video file
				f, err := os.Open("../samples/1280x720_2mb.mp4")
				require.NoError(t, err)
				defer f.Close()

				// Post video upload
				_, body, _ := session.PostMultipart(pathUpload, videoTitle, "video.avi", f)

				// Reading the body
				rawBody, _ := ioutil.ReadAll(body)
				var uploadResponse helpers.Response
				_ = json.Unmarshal(rawBody, &uploadResponse)

				// Update videoID and Get video status
				videoID = uploadResponse.Video.ID
			})

			g.It("Delete a video", func() {
				code, err := session.Delete("/api/v1/videos/" + videoID + "/delete")
				require.NoError(t, err)
				require.Equal(t, 200, code)
			})
		})
	})
}
