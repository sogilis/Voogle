package encoding

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/pkg/clients"
	contracts "github.com/Sogilis/Voogle/src/pkg/contracts/v1"
	"github.com/Sogilis/Voogle/src/pkg/ffmpeg"
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
	withSound, err := ffmpeg.CheckContainsSound(data.GetSource())
	if err != nil {
		return err
	}
	if !withSound {
		if err := ffmpeg.AddEmptyAudioTrack(data.GetSource()); err != nil {
			return err
		}
	}

	res, err := ffmpeg.ExtractResolution(data.GetSource())
	if err != nil {
		return err
	}
	if err = ffmpeg.ConvertToHLS(data.GetSource(), res); err != nil {
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
