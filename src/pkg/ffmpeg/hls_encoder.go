package ffmpeg

import (
	"fmt"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func ConvertToHLS(source string, res resolution) error {
	cmd, args, err := generateCommand(source, res)
	if err != nil {
		return err
	}

	log.Debug("FFMPEG command: ", cmd, strings.Join(args, " "))
	rawOutput, err := exec.Command(cmd, args...).CombinedOutput()
	log.Debug("FFMPEG output: ", string(rawOutput[:]))
	return err
}

func generateCommand(filepath string, res resolution) (string, []string, error) {
	// Example of the biggest command that can be generated
	// ffmpeg -y -i <filepath> \
	//              -pix_fmt yuv420p \
	//              -vcodec libx264 \
	//              -preset fast \
	//              -g 48 -sc_threshold 0 \
	//              -map 0:0 -map 0:1 -map 0:0 -map 0:1 -map 0:0 -map 0:1 -map 0:0 -map 0:1 \
	//              -s:v:0 640x480 -c:v:0 libx264 -b:v:0 1000k \
	//              -s:v:1 1280x720 -c:v:1 libx264 -b:v:1 2000k  \
	//              -s:v:2 1920x1080 -c:v:2 libx264 -b:v:2 4000k  \
	//              -s:v:3 3840x2160 -c:v:3 libx264 -b:v:3 8000k  \
	//              -c:a aac -b:a 128k -ac 2 \
	//              -var_stream_map "v:0,a:0 v:1,a:1 v:2,a:2 v:3,a:3" \
	//              -master_pl_name master.m3u8 \
	//              -f hls -hls_time 6 -hls_list_size 0 \
	//              -hls_segment_filename "v%v/segment%d.ts" \
	//              v%v/segment_index.m3u8

	if res.x < 640 && res.y < 480 {
		return "", nil, fmt.Errorf("resolution (%d,%d) is below minimal resolution (640x480)", res.x, res.y)
	}

	command := "ffmpeg"
	args := []string{"-y", "-i", filepath, "-pix_fmt", "yuv420p", "-vcodec", "libx264", "-preset", "fast", "-g", "48", "-sc_threshold", "0"}
	sound := []string{"-map", "0:0", "-map", "0:1"}
	resolutionTarget := []string{"-s:v:0", "640x480", "-c:v:0", "libx264", "-b:v:0", "1000k"}
	streamMap := "v:0,a:0"

	if res.GreaterOrEqualResolution(resolution{1280, 720}) {
		sound = append(sound, "-map", "0:0", "-map", "0:1")
		resolutionTarget = append(resolutionTarget, "-s:v:1", "1280x720", "-c:v:1", "libx264", "-b:v:1", "2000k")
		streamMap = streamMap + " v:1,a:1"
	}
	if res.GreaterOrEqualResolution(resolution{1920, 1080}) {
		sound = append(sound, "-map", "0:0", "-map", "0:1")
		resolutionTarget = append(resolutionTarget, "-s:v:2", "1920x1080", "-c:v:2", "libx264", "-b:v:2", "4000k")
		streamMap = streamMap + " v:2,a:2"
	}
	if res.GreaterOrEqualResolution(resolution{3840, 2160}) {
		sound = append(sound, "-map", "0:0", "-map", "0:1")
		resolutionTarget = append(resolutionTarget, "-s:v:3", "3840x2160", "-c:v:3", "libx264", "-b:v:3", "8000k")
		streamMap = streamMap + " v:3,a:3"
	}

	args = append(args, sound...)
	args = append(args, resolutionTarget...)
	args = append(args, "-c:a", "aac", "-b:a", "128k", "-ac", "2")
	args = append(args, "-var_stream_map", streamMap)
	args = append(args, "-master_pl_name", "master.m3u8", "-f", "hls", "-hls_time", "2", "-hls_list_size", "0", "-hls_segment_filename", "v%v/segment%d.ts", "v%v/segment_index.m3u8")

	return command, args, nil
}
