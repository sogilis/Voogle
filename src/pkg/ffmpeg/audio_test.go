package ffmpeg_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/Sogilis/Voogle/src/pkg/ffmpeg"
)

func Test_videoAddSound(t *testing.T) {
	t.SkipNow()
	cases := []struct {
		Name          string
		GivenFilepath string
		ExpectSound   bool
		ExpectError   bool
	}{
		{Name: "Without Sound", GivenFilepath: "../../../samples/video_without_sound.mp4", ExpectSound: false, ExpectError: false},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			sound, err := CheckContainsSound(tt.GivenFilepath)
			require.NoError(t, err)
			require.Equal(t, tt.ExpectSound, sound)

			// Save input video in order to restore it at the end of the test
			inputVideo, err := os.ReadFile(tt.GivenFilepath)
			require.NoError(t, err)

			err = AddEmptyAudioTrack(tt.GivenFilepath)
			require.NoError(t, err)
			defer func() {
				err := os.WriteFile(tt.GivenFilepath, inputVideo, 0666)
				require.NoError(t, err)
			}()

			sound, err = CheckContainsSound(tt.GivenFilepath)
			require.NoError(t, err)
			require.Equal(t, true, sound)
		})
	}
}
