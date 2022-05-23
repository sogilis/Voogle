package ffmpeg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ExtractResolution(t *testing.T) {
	//  LFS Github quota is not really fair, is uses bandwidth even when we use it within the Github Actions CI
	//  So until we've found an alternative, we won't test the video processing part
	t.SkipNow()
	cases := []struct {
		GivenPath        string
		GivenFilename    string
		ExpectResolution resolution
		ExpectError      bool
	}{
		{
			GivenPath:        "../../../../samples/", // FIXME(JPR): Root of the project from the test file (We need may need a better way to address these)
			GivenFilename:    "320x240_testvideo.mp4",
			ExpectResolution: resolution{320, 240},
			ExpectError:      false,
		},
		{
			GivenPath:        "../../../../samples/", // FIXME(JPR): Root of the project from the test file (We need may need a better way to address these)
			GivenFilename:    "960x400_ocean_with_audio.avi",
			ExpectResolution: resolution{960, 400},
			ExpectError:      false,
		},
		{
			GivenPath:        "../../../../samples/", // FIXME(JPR): Root of the project from the test file (We need may need a better way to address these)
			GivenFilename:    "4K-10bit.mkv",
			ExpectResolution: resolution{3840, 2160},
			ExpectError:      false,
		},
		{
			GivenPath:        "../../../../samples/", // FIXME(JPR): Root of the project from the test file (We need may need a better way to address these)
			GivenFilename:    "960x400_ocean_with_audio.mkv",
			ExpectResolution: resolution{960, 400},
			ExpectError:      false,
		},
		{
			GivenPath:        "../../../../samples/", // FIXME(JPR): Root of the project from the test file (We need may need a better way to address these)
			GivenFilename:    "1280x720_2mb.mp4",
			ExpectResolution: resolution{1280, 720},
			ExpectError:      false,
		},
	}

	for _, tt := range cases {
		t.Run("Extract resolution from video "+tt.GivenFilename, func(t *testing.T) {
			res, err := ExtractResolution(tt.GivenPath + tt.GivenFilename)
			if tt.ExpectError {
				require.NotNil(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, res.x, tt.ExpectResolution.x)
			require.Equal(t, res.y, tt.ExpectResolution.y)
		})
	}
}

func Test_videoHaveSound(t *testing.T) {
	t.SkipNow()
	cases := []struct {
		Name          string
		GivenFilepath string
		ExpectSound   bool
		ExpectError   bool
	}{
		{Name: "With Sound", GivenFilepath: "../../../samples/1280x720_2mb.mp4", ExpectSound: true, ExpectError: false},
		{Name: "Without Sound", GivenFilepath: "../../../samples/video_without_sound.mp4", ExpectSound: false, ExpectError: false},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			sound, err := CheckContainsSound(tt.GivenFilepath)
			if tt.ExpectError {
				require.NotNil(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.ExpectSound, sound)
		})
	}
}
