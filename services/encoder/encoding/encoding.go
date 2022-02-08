package encoding

import (
	"context"
	"fmt"
	contracts "github.com/Sogilis/Voogle/services/encoder/contracts/v1"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/services/encoder/clients"
)

// Process input video into a HLS video
func Process(s3Client clients.IS3Client, data *contracts.Video) error {
	// Going to the working directory
	processingFolder := filepath.Join(os.TempDir(), "/encoder-processing-dir")
	if err := os.MkdirAll(processingFolder, os.ModePerm); err != nil {
		return err
	}
	if err := os.Chdir(processingFolder); err != nil {
		return err
	}

	// Download and write the source file on the filesystem
	source, err := s3Client.GetObject(context.Background(), data.GetId()+"/"+data.GetSource())
	if err != nil {
		return err
	}
	f, err := os.Create(data.GetSource())
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, source); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	// Video processing
	res, err := extractResolution(data.GetSource())
	if err != nil {
		return err
	}
	command, args, err := generateCommand(data.GetSource(), res)
	if err != nil {
		return err
	}
	if err = convertToHLS(command, args); err != nil {
		return err
	}

	log.Info("Processing of video ", data.GetId(), "done - Uploading to S3")
	// Uploading files to the S3
	err = filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path == "." || (!strings.HasSuffix(path, ".ts") && !strings.HasSuffix(path, ".m3u8")) {
				log.Debug("Skipping ", path)
				return nil
			}
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer func() { _ = f.Close() }()
			return s3Client.PutObjectInput(context.Background(), f, filepath.Join(data.GetId(), path))
		})
	if err != nil {
		return err
	}
	if err := s3Client.PutObjectInput(context.Background(), strings.NewReader(""), data.GetId()+"/Ready.txt"); err != nil {
		return err
	}

	// CLeaning up
	if err = os.Chdir(os.TempDir()); err != nil {
		return err
	}
	return os.RemoveAll(processingFolder)
}

type resolution struct {
	x uint64
	y uint64
}

func (r resolution) GreaterOrEqualResolution(input resolution) bool {
	return r.x >= input.x && r.y >= input.y
}

// Extract resolution of the video
func extractResolution(filepath string) (resolution, error) {
	// ffprobe -v error -select_streams v:0 -show_entries stream=width,height -of csv=s=x:p=0 <filepath>
	rawOutput, err := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=s=x:p=0", filepath).Output()
	if err != nil {
		return resolution{}, err
	}
	output := string(rawOutput[:])

	//Sometimes, ffprobe return several resolution despite the video only have one video track
	firstLine := strings.Split(output, "\n")[0] // We get: XRESxYRES

	splitResolution := strings.Split(firstLine, "x")
	var x, y uint64
	if x, err = strconv.ParseUint(splitResolution[0], 10, 32); err != nil {
		return resolution{}, err
	}
	if y, err = strconv.ParseUint(splitResolution[1], 10, 32); err != nil {
		return resolution{}, err
	}

	return resolution{x, y}, nil
}

func convertToHLS(cmd string, args []string) error {
	log.Debug("FFMPEG command: ", cmd, strings.Join(args, " "))
	rawOutput, err := exec.Command(cmd, args...).CombinedOutput()
	log.Debug("FFMPEG output: ", string(rawOutput[:]))
	return err
}

func generateCommand(filepath string, res resolution) (string, []string, error) {
	// ffmpeg -y -i <filepath> \
	//              -pix_fmt yuv420p \
	//              -vcodec libx264 \
	//              -preset slow \
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
	args = append(args, "-master_pl_name", "master.m3u8", "-f", "hls", "-hls_time", "6", "-hls_list_size", "0", "-hls_segment_filename", "v%v/segment%d.ts", "v%v/segment_index.m3u8")

	return command, args, nil
}
