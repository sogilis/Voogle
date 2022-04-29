package end2end_tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/franela/goblin"
	"github.com/stretchr/testify/assert"

	"github.com/Sogilis/Voogle/end2end/helpers"
)

func Test_Videos(t *testing.T) {
	host := os.Getenv("E2E_API_ENDPOINT")

	user := os.Getenv("E2E_USER_NAME")
	pwd := os.Getenv("E2E_USER_PWD")

	fmt.Println("'", user, "'")
	fmt.Println("'", pwd, "'")

	g := Goblin(t)
	g.Describe("Videos >", func() {
		g.Before(func() {
			// TODO clean DATA
		})

		g.Describe("List >", func() {
			path := "/api/v1/videos/list"

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
					var videoData helpers.AllVideos
					err = json.Unmarshal(rawBody, &videoData)
					assert.NoError(t, err)

					assert.Equal(t, "Success", videoData.Status)
					assert.Equal(t, 0, len(videoData.Data))
				})

				g.It("Returns a list of videos with One element", func() {
					t.Log("PATH - GET - " + path)

					f, err := os.Open("../samples/1280x720_2mb.mp4")
					assert.NoError(t, err)

					code, _, err := session.PostMultipart("/api/v1/videos/upload", "test data", "video.avi", f)
					assert.NoError(t, err)
					assert.Equal(t, 200, code)

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
				})
			})
		})
	})
}
