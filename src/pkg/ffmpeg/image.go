package ffmpeg

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func ConvertImg(srcpath string, dstpath string) error {
	_, err := exec.Command("ffmpeg", "-i", srcpath, dstpath).CombinedOutput()
	if err != nil {
		return err
	}
	log.Debug("Image converted")
	return nil
}
