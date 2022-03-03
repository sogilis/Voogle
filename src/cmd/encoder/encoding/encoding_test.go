package encoding

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
			res, err := extractResolution(tt.GivenPath + tt.GivenFilename)
			if tt.ExpectError {
				assert.NotNil(t, err)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, res.x, tt.ExpectResolution.x)
			assert.Equal(t, res.y, tt.ExpectResolution.y)
		})
	}
}

func Test_GenerateCommand(t *testing.T) {
	cases := []struct {
		Name            string
		GivenFilePath   string
		GivenResolution resolution
		ExpectCommand   string
		ExpectArgs      string
		ExpectError     bool
	}{
		{
			Name:            "Resolution below minimal",
			GivenFilePath:   "someName.mp4",
			GivenResolution: resolution{x: 0, y: 0},
			ExpectCommand:   "",
			ExpectError:     true,
		},
		{
			Name:            "With resolution: 640x480",
			GivenFilePath:   "someName.mp4",
			GivenResolution: resolution{x: 640, y: 480},
			ExpectCommand:   "ffmpeg",
			ExpectArgs:      "-y -i someName.mp4 -pix_fmt yuv420p -vcodec libx264 -preset fast -g 48 -sc_threshold 0 -map 0:0 -map 0:1 -s:v:0 640x480 -c:v:0 libx264 -b:v:0 1000k -c:a aac -b:a 128k -ac 2 -var_stream_map v:0,a:0 -master_pl_name master.m3u8 -f hls -hls_time 6 -hls_list_size 0 -hls_segment_filename v%v/segment%d.ts v%v/segment_index.m3u8",
			ExpectError:     false,
		},
		{
			Name:            "With resolution: 1280x720",
			GivenFilePath:   "someName.mp4",
			GivenResolution: resolution{x: 1280, y: 720},
			ExpectCommand:   "ffmpeg",
			ExpectArgs:      "-y -i someName.mp4 -pix_fmt yuv420p -vcodec libx264 -preset fast -g 48 -sc_threshold 0 -map 0:0 -map 0:1 -map 0:0 -map 0:1 -s:v:0 640x480 -c:v:0 libx264 -b:v:0 1000k -s:v:1 1280x720 -c:v:1 libx264 -b:v:1 2000k -c:a aac -b:a 128k -ac 2 -var_stream_map v:0,a:0 v:1,a:1 -master_pl_name master.m3u8 -f hls -hls_time 6 -hls_list_size 0 -hls_segment_filename v%v/segment%d.ts v%v/segment_index.m3u8",
			ExpectError:     false,
		},
		{
			Name:            "With resolution 1920x1080",
			GivenFilePath:   "someName.mp4",
			GivenResolution: resolution{x: 1920, y: 1080},
			ExpectCommand:   "ffmpeg",
			ExpectArgs:      "-y -i someName.mp4 -pix_fmt yuv420p -vcodec libx264 -preset fast -g 48 -sc_threshold 0 -map 0:0 -map 0:1 -map 0:0 -map 0:1 -map 0:0 -map 0:1 -s:v:0 640x480 -c:v:0 libx264 -b:v:0 1000k -s:v:1 1280x720 -c:v:1 libx264 -b:v:1 2000k -s:v:2 1920x1080 -c:v:2 libx264 -b:v:2 4000k -c:a aac -b:a 128k -ac 2 -var_stream_map v:0,a:0 v:1,a:1 v:2,a:2 -master_pl_name master.m3u8 -f hls -hls_time 6 -hls_list_size 0 -hls_segment_filename v%v/segment%d.ts v%v/segment_index.m3u8",
			ExpectError:     false,
		},
		{
			Name:            "With resolution 3840x2160",
			GivenFilePath:   "someName.mp4",
			GivenResolution: resolution{x: 3840, y: 2160},
			ExpectCommand:   "ffmpeg",
			ExpectArgs:      "-y -i someName.mp4 -pix_fmt yuv420p -vcodec libx264 -preset fast -g 48 -sc_threshold 0 -map 0:0 -map 0:1 -map 0:0 -map 0:1 -map 0:0 -map 0:1 -map 0:0 -map 0:1 -s:v:0 640x480 -c:v:0 libx264 -b:v:0 1000k -s:v:1 1280x720 -c:v:1 libx264 -b:v:1 2000k -s:v:2 1920x1080 -c:v:2 libx264 -b:v:2 4000k -s:v:3 3840x2160 -c:v:3 libx264 -b:v:3 8000k -c:a aac -b:a 128k -ac 2 -var_stream_map v:0,a:0 v:1,a:1 v:2,a:2 v:3,a:3 -master_pl_name master.m3u8 -f hls -hls_time 6 -hls_list_size 0 -hls_segment_filename v%v/segment%d.ts v%v/segment_index.m3u8",
			ExpectError:     false,
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			cmd, args, err := generateCommand(tt.GivenFilePath, tt.GivenResolution)
			if tt.ExpectError {
				assert.NotNil(t, err)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, tt.ExpectCommand, cmd)
			assert.Equal(t, tt.ExpectArgs, strings.Join(args, " "))
		})
	}
}

func Test_convertToHLS(t *testing.T) {
	// This test is really CPU intensive and takes time
	t.SkipNow()
	cases := []struct {
		Name            string
		GivenFilePath   string
		GivenResolution resolution
		ExpectError     bool
	}{
		{Name: "Low quality video (960x400_ocean_with_audio.avi)", GivenFilePath: "../../../../../samples/960x400_ocean_with_audio.avi", GivenResolution: resolution{960, 400}, ExpectError: false},
		{Name: "Medium low quality video (1280x720_2mb.mp4)", GivenFilePath: "../../../../../samples/1280x720_2mb.mp4", GivenResolution: resolution{1280, 720}, ExpectError: false},
		{Name: "High quality video (4K-10bit.mkv)", GivenFilePath: "../../../../../samples/4K-10bit.mkv", GivenResolution: resolution{3840, 2160}, ExpectError: false},
		{Name: "Video that doesn't exists", GivenFilePath: "../../../../../samples/none.mkv", GivenResolution: resolution{3840, 2160}, ExpectError: true},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			_ = os.Mkdir("tmpVideoTest", os.ModePerm)
			_ = os.Chdir("tmpVideoTest")
			cmd, args, err := generateCommand(tt.GivenFilePath, tt.GivenResolution)
			assert.NoError(t, err)
			err = convertToHLS(cmd, args)
			if tt.ExpectError {
				assert.NotNil(t, err)
				return
			}
			assert.NoError(t, err)
			_ = os.Chdir("..")
			_ = os.RemoveAll("tmpVideoTest")
		})
	}
}
