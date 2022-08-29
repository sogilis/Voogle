package ffmpeg

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func ConvertImg(srcpath string, dstpath string) error {
	_, err := exec.Command("ffmpeg", "-i", srcpath, "-vf",
		"crop='if(gt(iw, ih), ih, iw)':'if(gt(iw, ih), ih, iw )':'if(gt(iw, ih), (iw-ih)/2, 0)':'if(gt(iw,ih), 0, (ih-iw)/2)', scale=250:250",
		dstpath).CombinedOutput()
	if err != nil {
		return err
	}
	log.Debug("Image converted")
	return nil
}
