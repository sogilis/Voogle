package end2end_tests

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Sogilis/Voogle/end2end/helpers"
	. "github.com/franela/goblin"
	"github.com/stretchr/testify/require"
)

// TODO : These tests are not independent of each other.
// To do so, we need a endpoint to delete video from the database.
func Test_Videos(t *testing.T) {
	host := os.Getenv("E2E_API_ENDPOINT")

	user := os.Getenv("E2E_USER_NAME")
	pwd := os.Getenv("E2E_USER_PWD")

	g := Goblin(t)
	g.Describe("Videos >", func() {
		// g.Describe("Upload >", func() {
		// 	g.Before(func() {
		// 		// TODO clean DATA
		// 	})

		// 	path := "/api/v1/videos/upload"
		// 	videoTitle := "test upload"

		// 	session := helpers.NewSession(host)
		// 	require.Nil(t, session.Login(user, pwd))

		// 	g.It("Upload one video", func() {
		// 		t.Log("PATH - POST - " + path)

		// 		f, err := os.Open("../samples/1280x720_2mb.mp4")
		// 		require.NoError(t, err)
		// 		defer f.Close()

		// 		// Post video upload
		// 		code, body, err := session.PostMultipart("/api/v1/videos/upload", videoTitle, "video.avi", f)
		// 		require.NoError(t, err)

		// 		require.Equal(t, 200, code)

		// 		// Reading the body
		// 		rawBody, err := ioutil.ReadAll(body)
		// 		require.NoError(t, err)
		// 		var uploadResponse helpers.Response
		// 		err = json.Unmarshal(rawBody, &uploadResponse)
		// 		require.NoError(t, err)
		// 		require.Equal(t, uploadResponse.Video.Title, videoTitle)
		// 	})

		// 	g.It("Returns an error title already exist", func() {
		// 		t.Log("PATH - POST - " + path)

		// 		f, err := os.Open("../samples/1280x720_2mb.mp4")
		// 		require.NoError(t, err)
		// 		defer f.Close()

		// 		// Post video upload
		// 		code, _, err := session.PostMultipart("/api/v1/videos/upload", videoTitle, "video.avi", f)
		// 		require.NoError(t, err)

		// 		require.Equal(t, 409, code)
		// 	})

		// 	g.It("Returns an error unsported media format", func() {
		// 		t.Log("PATH - POST - " + path)
		// 		imageTitle := "test upload image"
		// 		f, err := os.Open("../samples/image.mp4")
		// 		require.NoError(t, err)
		// 		defer f.Close()

		// 		// Post video upload
		// 		code, _, err := session.PostMultipart("/api/v1/videos/upload", imageTitle, "falsevideo.mp4", f)
		// 		require.NoError(t, err)

		// 		require.Equal(t, 415, code)
		// 	})
		// })

		// g.Describe("List >", func() {
		// 	path := "/api/v1/videos/list"

		// 	g.Describe("Without login >", func() {
		// 		session := helpers.NewSession(host)
		// 		g.It("Returns a 401", func() {
		// 			t.Log("PATH - GET - " + path)

		// 			code, _, err := session.Get(path)
		// 			require.NoError(t, err)

		// 			require.Equal(t, 401, code)
		// 		})
		// 	})

		// 	g.Describe("With login >", func() {
		// 		session := helpers.NewSession(host)
		// 		require.Nil(t, session.Login(user, pwd))

		// 		g.Describe("With empty list >", func() {
		// 			g.It("Returns an empty list of videos", func() {
		// 				t.Log("PATH - GET - " + path)

		// 				// Get video list
		// 				code, body, err := session.Get(path)
		// 				require.NoError(t, err)

		// 				require.Equal(t, 200, code)

		// 				// Reading the body
		// 				rawBody, err := ioutil.ReadAll(body)
		// 				require.NoError(t, err)
		// 				var videoData helpers.AllVideos
		// 				err = json.Unmarshal(rawBody, &videoData)
		// 				require.NoError(t, err)

		// 				require.Equal(t, "Success", videoData.Status)
		// 				require.Equal(t, 0, len(videoData.Data))
		// 			})
		// 		})

		// 		g.Describe("With one video >", func() {

		// 			videoTitle := "test list"
		// 			g.Before(func() {
		// 				// TODO clean DATA

		// 				path := "/api/v1/videos/upload"

		// 				session := helpers.NewSession(host)
		// 				require.Nil(t, session.Login(user, pwd))

		// 				t.Log("PATH - POST - " + path)

		// 				f, err := os.Open("../samples/1280x720_2mb.mp4")
		// 				require.NoError(t, err)
		// 				defer f.Close()

		// 				// Post video upload
		// 				code, body, err := session.PostMultipart(path, videoTitle, "video.avi", f)
		// 				require.NoError(t, err)

		// 				require.Equal(t, 200, code)

		// 				// Reading the body
		// 				rawBody, err := ioutil.ReadAll(body)
		// 				require.NoError(t, err)
		// 				var uploadResponse helpers.Response
		// 				err = json.Unmarshal(rawBody, &uploadResponse)
		// 				require.NoError(t, err)

		// 				// Get video status
		// 				g.Timeout(time.Duration(60) * time.Second)
		// 				err = session.WaitVideoEncoded("/api/v1/videos/" + uploadResponse.Video.ID + "/status")
		// 				require.NoError(t, err)

		// 			})

		// 			g.It("Returns a list of videos with One element", func() {
		// 				t.Log("PATH - GET - " + path)

		// 				// Get video list
		// 				code, body, err := session.Get(path)
		// 				require.NoError(t, err)

		// 				require.Equal(t, 200, code)

		// 				// Reading the body
		// 				rawBody, err := ioutil.ReadAll(body)
		// 				require.NoError(t, err)
		// 				var videoData helpers.AllVideos
		// 				err = json.Unmarshal(rawBody, &videoData)
		// 				require.NoError(t, err)

		// 				require.Equal(t, "Success", videoData.Status)
		// 				require.Equal(t, 1, len(videoData.Data))
		// 				require.Equal(t, videoTitle, videoData.Data[0].Title)
		// 			})
		// 		})
		// 	})
		// })

		g.Describe("Stream >", func() {
			videoTitle := "test stream"

			g.Before(func() {
				// TODO clean DATA

				path := "/api/v1/videos/upload"

				session := helpers.NewSession(host)
				require.Nil(t, session.Login(user, pwd))

				t.Log("PATH - POST - " + path)

				f, err := os.Open("../samples/1280x720_2mb.mp4")
				require.NoError(t, err)
				defer f.Close()

				// Post video upload
				code, body, err := session.PostMultipart(path, videoTitle, "video.avi", f)
				require.NoError(t, err)

				require.Equal(t, 200, code)

				// Reading the body
				rawBody, err := ioutil.ReadAll(body)
				require.NoError(t, err)
				var uploadResponse helpers.Response
				err = json.Unmarshal(rawBody, &uploadResponse)
				require.NoError(t, err)

				// Get video status
				g.Timeout(60 * time.Second)
				err = session.WaitVideoEncoded("/api/v1/videos/" + uploadResponse.Video.ID + "/status")
				require.NoError(t, err)

			})

			g.Describe("With login >", func() {

				session := helpers.NewSession(host)
				require.Nil(t, session.Login(user, pwd))

				g.It("Returns video stream master and first part", func() {
					// Get video list
					code, body, err := session.Get("/api/v1/videos/list")
					require.NoError(t, err)

					require.Equal(t, 200, code)

					// Reading the body
					rawBody, err := ioutil.ReadAll(body)
					require.NoError(t, err)
					var videoData helpers.AllVideos
					err = json.Unmarshal(rawBody, &videoData)
					require.NoError(t, err)

					require.Equal(t, "Success", videoData.Status)
					require.Equal(t, 1, len(videoData.Data))
					require.Equal(t, videoTitle, videoData.Data[0].Title)

					// Get video master
					code, body, err = session.Get("/api/v1/videos/" + videoData.Data[0].Id + "/streams/master.m3u8")
					require.NoError(t, err)

					require.Equal(t, 200, code)

					// Reading the body
					rawBody, err = ioutil.ReadAll(body)
					require.NoError(t, err)

					//Parse response
					parseBody := strings.Split(string(rawBody), "#")
					require.Equal(t, "EXTM3U\n", parseBody[1])
					require.Equal(t, "EXT-X-VERSION:3\n", parseBody[2])

					// Get video part
					code, _, err = session.Get("/api/v1/videos/" + videoData.Data[0].Id + "/streams/v0/segment0.ts")
					require.NoError(t, err)

					require.Equal(t, 200, code)
				})
			})
		})
	})
}
