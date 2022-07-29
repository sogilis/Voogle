package integration_tests

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/franela/goblin"
	"github.com/stretchr/testify/require"

	"github.com/Sogilis/Voogle/integration/helpers"
)

const GOBLIN_TEST_TIMEOUT time.Duration = time.Second * 15

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

		videoLocation := "../samples/1280x720_2mb.mp4"
		videoTitle := "test"
		var videoID string

		jpgCoverLocation := "../samples/cover.jpg"
		pngCoverLocation := "../samples/cover.png"

		// Open first HLS video part
		videoPart, err := os.Open("../samples/1280x720_2mb_segment0.ts")
		require.NoError(t, err)
		defer videoPart.Close()

		g.AfterEach(func() {
			// Clear data
			_, _, _ = session.Put("/api/v1/videos/" + videoID + "/archive")
			_, _ = session.Delete("/api/v1/videos/" + videoID + "/delete")
		})

		//////////////////
		// UPLOAD VIDEO //
		//////////////////
		g.Describe("Upload >", func() {
			g.Describe("With video and no cover image >", func() {
				g.It("Upload one video", func() {
					t.Log("PATH - POST - " + pathUpload)

					g.Timeout(GOBLIN_TEST_TIMEOUT)

					// Open video file
					f, err := os.Open(videoLocation)
					require.NoError(t, err)
					defer f.Close()

					// Post video upload
					code, body, err := session.PostMultipart(pathUpload, videoTitle, "video.avi", f, nil)
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
					_ = session.WaitVideoEncoded("/api/v1/videos/" + videoID + "/status")
				})
			})

			g.Describe("With video and jpg cover image >", func() {
				g.It("Upload one video and jpg cover", func() {
					t.Log("PATH - POST - " + pathUpload)

					g.Timeout(GOBLIN_TEST_TIMEOUT)

					// Open video file
					fVideo, err := os.Open(videoLocation)
					require.NoError(t, err)
					defer fVideo.Close()

					// Open cover file
					fCover, err := os.Open(jpgCoverLocation)
					require.NoError(t, err)
					defer fCover.Close()

					// Post video upload
					code, body, err := session.PostMultipart(pathUpload, videoTitle, "video.avi", fVideo, fCover)
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
					_ = session.WaitVideoEncoded("/api/v1/videos/" + videoID + "/status")
				})
			})

			g.Describe("With video and png cover image >", func() {
				g.It("Upload one video and png cover", func() {
					t.Log("PATH - POST - " + pathUpload)

					g.Timeout(GOBLIN_TEST_TIMEOUT)

					// Open video file
					fVideo, err := os.Open(videoLocation)
					require.NoError(t, err)
					defer fVideo.Close()

					// Open cover file
					fCover, err := os.Open(pngCoverLocation)
					require.NoError(t, err)
					defer fCover.Close()

					// Post video upload
					code, body, err := session.PostMultipart(pathUpload, videoTitle, "video.avi", fVideo, fCover)
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
					_ = session.WaitVideoEncoded("/api/v1/videos/" + videoID + "/status")
				})
			})

			g.Describe("With video title already exists >", func() {
				g.Before(func() {
					uploadVideoWaitForEncode(&videoLocation, &pathUpload, &videoTitle, &videoID, session)
				})

				g.It("Returns an error title already exist", func() {
					t.Log("PATH - POST - " + pathUpload)

					g.Timeout(GOBLIN_TEST_TIMEOUT)

					// Open video file
					f, err := os.Open(videoLocation)
					require.NoError(t, err)
					defer f.Close()

					// Post video upload with same title
					code, _, err := session.PostMultipart(pathUpload, videoTitle, "video.avi", f, nil)
					require.NoError(t, err)
					require.Equal(t, 409, code)
				})

			})

			g.Describe("With image as video file >", func() {
				g.It("Returns an error unsported media format", func() {
					t.Log("PATH - POST - " + pathUpload)

					g.Timeout(GOBLIN_TEST_TIMEOUT)

					// Open image file
					fImage, err := os.Open("../samples/image.mp4")
					require.NoError(t, err)
					defer fImage.Close()

					// Post video upload with image file
					code, _, err := session.PostMultipart(pathUpload, "test image", "falsevideo.mp4", fImage, nil)
					require.NoError(t, err)
					require.Equal(t, 415, code)
				})
			})
		})

		////////////////
		// LIST VIDEO //
		////////////////
		g.Describe("List >", func() {
			g.Describe("Without login >", func() {
				g.It("Returns a 401", func() {
					t.Log("PATH - GET - " + pathList)

					g.Timeout(GOBLIN_TEST_TIMEOUT)

					code, _, err := sessionNil.Get(pathList)
					require.NoError(t, err)
					require.Equal(t, 401, code)
				})
			})

			g.Describe("With login >", func() {
				g.Describe("With empty list >", func() {
					g.It("Returns an empty list of videos", func() {
						t.Log("PATH - GET - " + pathList)

						g.Timeout(GOBLIN_TEST_TIMEOUT)

						// Get video list
						code, body, err := session.Get(pathList)
						require.NoError(t, err)
						require.Equal(t, 200, code)

						// Reading the body
						rawBody, err := ioutil.ReadAll(body)
						require.NoError(t, err)

						// Retrieve videos list informations
						var videoData helpers.VideoListResponse
						err = json.Unmarshal(rawBody, &videoData)
						require.NoError(t, err)

						require.Equal(t, 0, len(videoData.Videos))
					})
				})

				g.Describe("With one video >", func() {
					g.Before(func() {
						uploadVideoWaitForEncode(&videoLocation, &pathUpload, &videoTitle, &videoID, session)
					})

					g.It("Returns a list of videos with One element", func() {
						t.Log("PATH - GET - " + pathList)

						g.Timeout(GOBLIN_TEST_TIMEOUT)

						// Get video list
						code, body, err := session.Get(pathList)
						require.NoError(t, err)
						require.Equal(t, 200, code)

						// Reading the body
						rawBody, err := ioutil.ReadAll(body)
						require.NoError(t, err)

						// Retrieve videos list informations
						var videoData helpers.VideoListResponse
						err = json.Unmarshal(rawBody, &videoData)
						require.NoError(t, err)

						require.Equal(t, 1, len(videoData.Videos))
						require.Equal(t, videoTitle, videoData.Videos[0].Title)
					})
				})
			})
		})

		//////////////////
		// STATUS VIDEO //
		//////////////////
		g.Describe("Status >", func() {
			g.Before(func() {
				uploadVideoWaitForEncode(&videoLocation, &pathUpload, &videoTitle, &videoID, session)
			})

			g.It("Get video status complete", func() {

				g.Timeout(GOBLIN_TEST_TIMEOUT)

				code, body, err := session.Get("/api/v1/videos/" + videoID + "/status")
				require.NoError(t, err)
				require.Equal(t, 200, code)

				// Reading the body
				rawBody, err := ioutil.ReadAll(body)
				require.NoError(t, err)

				// Retrieve video status
				var statusResponse helpers.VideoStatus
				err = json.Unmarshal(rawBody, &statusResponse)
				require.NoError(t, err)
				require.Equal(t, statusResponse.Title, videoTitle)
				require.Equal(t, strings.ToLower(statusResponse.Status), "complete")
			})
		})

		////////////////
		// INFO VIDEO //
		////////////////
		g.Describe("Info >", func() {
			g.Before(func() {
				uploadVideoWaitForEncode(&videoLocation, &pathUpload, &videoTitle, &videoID, session)
			})

			g.It("Get video infos", func() {

				g.Timeout(GOBLIN_TEST_TIMEOUT)

				code, body, err := session.Get("/api/v1/videos/" + videoID + "/info")
				require.NoError(t, err)
				require.Equal(t, 200, code)

				// Reading the body
				rawBody, err := ioutil.ReadAll(body)
				require.NoError(t, err)

				// Retrieve video informations
				var statusResponse helpers.VideoInfo
				err = json.Unmarshal(rawBody, &statusResponse)
				require.NoError(t, err)
				require.Equal(t, statusResponse.Title, videoTitle)
			})
		})

		//////////////////
		// STREAM VIDEO //
		//////////////////
		g.Describe("Stream >", func() {
			g.Describe("Get video master >", func() {
				g.Before(func() {
					uploadVideoWaitForEncode(&videoLocation, &pathUpload, &videoTitle, &videoID, session)
				})

				g.It("Returns video stream master", func() {

					g.Timeout(GOBLIN_TEST_TIMEOUT)

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
				})
			})

			g.Describe("Get first video part without transformation >", func() {
				g.Before(func() {
					uploadVideoWaitForEncode(&videoLocation, &pathUpload, &videoTitle, &videoID, session)
				})

				g.It("Returns first video part", func() {

					g.Timeout(GOBLIN_TEST_TIMEOUT)

					// Get video part
					code, body, err := session.Get("/api/v1/videos/" + videoID + "/streams/v0/segment0.ts")
					require.NoError(t, err)
					require.NotEmpty(t, body)
					require.Equal(t, 200, code)
				})
			})

			g.Describe("Get first video part with gray transformation >", func() {
				g.Before(func() {
					uploadVideoWaitForEncode(&videoLocation, &pathUpload, &videoTitle, &videoID, session)
				})

				g.It("Returns first video part gray", func() {

					g.Timeout(GOBLIN_TEST_TIMEOUT)

					// Get gray video part
					code, body, err := session.Get("/api/v1/videos/" + videoID + "/streams/v0/segment0.ts?filter=gray")
					require.NoError(t, err)
					require.NotEmpty(t, body)
					require.Equal(t, 200, code)
				})
			})

			g.Describe("Get first video part with flip transformation >", func() {
				g.Before(func() {
					uploadVideoWaitForEncode(&videoLocation, &pathUpload, &videoTitle, &videoID, session)
				})

				g.It("Returns first video part flip", func() {

					g.Timeout(GOBLIN_TEST_TIMEOUT)

					// Get flip video part
					code, body, err := session.Get("/api/v1/videos/" + videoID + "/streams/v0/segment0.ts?filter=flip")
					require.NoError(t, err)
					require.NotEmpty(t, body)
					require.Equal(t, 200, code)
				})
			})

			g.Describe("Get first video part with gray and flip transformation >", func() {
				g.Before(func() {
					uploadVideoWaitForEncode(&videoLocation, &pathUpload, &videoTitle, &videoID, session)
				})

				g.It("Returns first video part", func() {

					g.Timeout(GOBLIN_TEST_TIMEOUT)

					// Get video part
					code, body, err := session.Get("/api/v1/videos/" + videoID + "/streams/v0/segment0.ts?filter=gray&filter=flip")
					require.NoError(t, err)
					require.NotEmpty(t, body)
					require.Equal(t, 200, code)
				})
			})
		})

		//////////////////
		// DELETE VIDEO //
		//////////////////
		g.Describe("Delete >", func() {
			g.Before(func() {
				uploadVideoWaitForEncode(&videoLocation, &pathUpload, &videoTitle, &videoID, session)
				_, _, _ = session.Put("/api/v1/videos/" + videoID + "/archive")
			})

			g.It("Delete a video", func() {

				g.Timeout(GOBLIN_TEST_TIMEOUT)

				code, err := session.Delete("/api/v1/videos/" + videoID + "/delete")
				require.NoError(t, err)
				require.Equal(t, 200, code)
			})
		})

		///////////////////
		// ARCHIVE VIDEO //
		///////////////////
		g.Describe("Archive >", func() {
			g.Before(func() {
				uploadVideoWaitForEncode(&videoLocation, &pathUpload, &videoTitle, &videoID, session)
			})

			g.It("Archive a video", func() {

				g.Timeout(GOBLIN_TEST_TIMEOUT)

				code, _, err := session.Put("/api/v1/videos/" + videoID + "/archive")
				require.NoError(t, err)
				require.Equal(t, 200, code)
			})
		})

		/////////////////////
		// UNARCHIVE VIDEO //
		/////////////////////
		g.Describe("Unarchive >", func() {
			g.Describe("Without archived video >", func() {
				g.Before(func() {
					uploadVideoWaitForEncode(&videoLocation, &pathUpload, &videoTitle, &videoID, session)
				})
				g.It("Unarchive video fails", func() {

					g.Timeout(GOBLIN_TEST_TIMEOUT)

					code, _, err := session.Put("/api/v1/videos/" + videoID + "/unarchive")
					require.NoError(t, err)
					require.Equal(t, 400, code)
				})
			})

			g.Describe("Without archived video >", func() {
				g.Before(func() {
					uploadVideoWaitForEncode(&videoLocation, &pathUpload, &videoTitle, &videoID, session)
					_, _, _ = session.Put("/api/v1/videos/" + videoID + "/archive")
				})
				g.It("Unarchive video", func() {

					g.Timeout(GOBLIN_TEST_TIMEOUT)

					code, _, err := session.Put("/api/v1/videos/" + videoID + "/unarchive")
					require.NoError(t, err)
					require.Equal(t, 200, code)
				})
			})
		})

		/////////////////////
		// GET VIDEO COVER //
		/////////////////////
		g.Describe("Cover >", func() {
			g.Before(func() {
				uploadVideoWaitForEncode(&videoLocation, &pathUpload, &videoTitle, &videoID, session)
			})

			g.It("Get video cover", func() {

				g.Timeout(GOBLIN_TEST_TIMEOUT)

				code, _, err := session.Get("/api/v1/videos/" + videoID + "/cover")
				require.NoError(t, err)
				require.Equal(t, 200, code)
			})
		})

		//////////////////////////
		// GET TRANSFROMER LIST //
		//////////////////////////
		g.Describe("Transformer >", func() {
			g.It("Get video cover", func() {

				g.Timeout(GOBLIN_TEST_TIMEOUT)

				code, body, err := session.Get("/api/v1/videos/transformer/list")
				require.NoError(t, err)
				require.Equal(t, 200, code)

				// Reading the body
				rawBody, err := ioutil.ReadAll(body)
				require.NoError(t, err)

				// Retrieve transformers
				var transformerList helpers.TransformerServiceListResponse
				err = json.Unmarshal(rawBody, &transformerList)
				require.NoError(t, err)
				require.True(t, transformerList.Services[0].Name == "gray" || transformerList.Services[0].Name == "flip")
				require.True(t, transformerList.Services[1].Name == "gray" || transformerList.Services[1].Name == "flip")
			})
		})
	})
}

func uploadVideoWaitForEncode(videoLocation, pathUpload, videoTitle, videoID *string, session helpers.Session) {
	// Open video file
	f, _ := os.Open(*videoLocation)
	defer f.Close()

	// Post video upload
	_, body, _ := session.PostMultipart(*pathUpload, *videoTitle, "video.avi", f, nil)

	// Reading the body
	rawBody, _ := ioutil.ReadAll(body)
	var uploadResponse helpers.Response
	_ = json.Unmarshal(rawBody, &uploadResponse)

	// Update videoID and Get video status
	*videoID = uploadResponse.Video.ID
	_ = session.WaitVideoEncoded("/api/v1/videos/" + *videoID + "/status")
}
