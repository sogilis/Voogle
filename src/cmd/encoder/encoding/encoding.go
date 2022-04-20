package encoding

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"
)

// Process input video into a HLS video
func Process(s3Client clients.IS3Client, videoData *contracts.Video) error {
	// Going to the working directory
	processingFolder := filepath.Join(os.TempDir(), "/encoder-processing-dir")
	if err := os.MkdirAll(processingFolder, os.ModePerm); err != nil {
		return err
	}
	if err := os.Chdir(processingFolder); err != nil {
		return err
	}
	defer func() {
		// CLeaning up
		_ = os.Chdir(os.TempDir())
		_ = os.RemoveAll(processingFolder)
	}()

	// Download and write the source file on the filesystem
	err := fetchVideoSource(s3Client, videoData)
	if err != nil {
		return err
	}

	// Video processing
	// Some video doesn't contains audio and HLS can't handle it, so we add an empty track
	err = encode(videoData)
	if err != nil {
		return err
	}

	log.Info("Processing of video ", videoData.GetId(), "done - Uploading to S3")
	// Uploading files to the S3
	err = uploadFiles(s3Client, videoData)
	if err != nil {
		return err
	}

	return nil
}

func fetchVideoSource(s3Client clients.IS3Client, videoData *contracts.Video) error {
	source, err := s3Client.GetObject(context.Background(), videoData.GetId()+"/"+videoData.GetSource())
	if err != nil {
		return err
	}
	f, err := os.Create(videoData.GetSource())
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, source); err != nil {
		return err
	}
	return f.Close()
}

func encode(data *contracts.Video) error {
	withSound, err := checkContainsSound(data.GetSource())
	if err != nil {
		return err
	}
	if !withSound {
		if err := addEmptyAudioTrack(data.GetSource()); err != nil {
			return err
		}
	}

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
	return nil
}

func uploadFiles(s3Client clients.IS3Client, data *contracts.Video) error {
	return filepath.Walk(".",
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
}

func addEmptyAudioTrack(fileName string) error {
	// ffmpeg -f lavfi -i anullsrc=channel_layout=stereo:sample_rate=44100 -i <filepath> -c:v copy -c:a aac -shortest <filepath>
	tmpPath := "tmp" + filepath.Ext(fileName)
	rawOutput, err := exec.Command("ffmpeg", "-y", "-f", "lavfi", "-i", "anullsrc=channel_layout=stereo:sample_rate=44100", "-i", fileName, "-c:v", "copy", "-c:a", "aac", "-shortest", tmpPath).CombinedOutput()
	if rawOutput != nil {
		log.Info("Adding Empty audio Track", string(rawOutput[:]))
	}
	if err != nil {
		return err
	}
	return os.Rename(tmpPath, fileName)
}

func convertToHLS(cmd string, args []string) error {
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
	args = append(args, "-master_pl_name", "master.m3u8", "-f", "hls", "-hls_time", "6", "-hls_list_size", "0", "-hls_segment_filename", "v%v/segment%d.ts", "v%v/segment_index.m3u8")

	return command, args, nil
}
