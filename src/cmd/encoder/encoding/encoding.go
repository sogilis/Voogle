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
		log.Error("Failed to fetch video source")
		return err
	}

	// Video processing
	// Some video doesn't contains audio and HLS can't handle it, so we add an empty track
	err = encode(videoData)
	if err != nil {
		log.Error("Failed to encode video")
		return err
	}

	// Download and write the cover file on the filesystem
	isCoverFetch, err := fetchCoverSource(s3Client, videoData)
	if err != nil {
		log.Error("Failed to fetch cover image")
		return err
	}

	// Cover image compression
	if isCoverFetch {
		if err = compressCover(); err != nil {
			log.Error("Failed to compress cover image")
			return err
		}
	}

	log.Info("Processing of video ", videoData.GetId(), "done - Uploading to S3")
	// Uploading files to the S3
	err = uploadFiles(s3Client, videoData)
	if err != nil {
		log.Error("Failde to upload video data to S3")
		return err
	}

	return nil
}

func fetchVideoSource(s3Client clients.IS3Client, videoData *contracts.Video) error {
	source, err := s3Client.GetObject(context.Background(), videoData.GetSource())
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Base(videoData.GetSource()))
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, source); err != nil {
		return err
	}
	return f.Close()
}

func encode(data *contracts.Video) error {
	sourcefile := filepath.Base(data.GetSource())

	withSound, err := ffmpeg.CheckContainsSound(sourcefile)
	if err != nil {
		return err
	}
	if !withSound {
		if err := ffmpeg.AddEmptyAudioTrack(sourcefile); err != nil {
			return err
		}
	}

	res, err := ffmpeg.ExtractResolution(sourcefile)
	if err != nil {
		return err
	}
	if err = ffmpeg.ConvertToHLS(sourcefile, res); err != nil {
		return err
	}
	return nil
}

func fetchCoverSource(s3Client clients.IS3Client, videoData *contracts.Video) (isFileFetch bool, err error) {
	if filepath.Ext(videoData.GetCoverPath()) != ".png" {
		return false, nil
	}

	source, err := s3Client.GetObject(context.Background(), videoData.GetCoverPath())
	if err != nil {
		return false, err
	}
	f, err := os.Create(filepath.Base(videoData.GetCoverPath()))
	if err != nil {
		return false, err
	}
	if _, err := io.Copy(f, source); err != nil {
		return true, err
	}
	return true, f.Close()
}

func compressCover() error {
	err := ffmpeg.ConvertImg("cover.png", "cover.jpg")
	if err != nil {
		return err
	}
	return nil
}

func uploadFiles(s3Client clients.IS3Client, data *contracts.Video) error {
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path == "." || (!strings.HasSuffix(path, ".ts") && !strings.HasSuffix(path, ".m3u8") && !strings.HasSuffix(path, ".jpg")) {
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

	if _, err = os.Stat("cover.jpg"); err == nil {
		err = s3Client.RemoveObject(context.Background(), data.GetCoverPath())
		if err != nil {
			return err
		}
	}

	return nil
}
