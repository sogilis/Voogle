package ffmpeg

import (
	"bytes"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func TransformFlip(inputFilename, outputFilename string) (*bytes.Buffer, error) {
	command := "ffmpeg"
	args := []string{"-i", inputFilename}
	args = append(args, "-muxdelay", "0", "-map", "0:0", "-map", "0:1", "-acodec", "copy", "-vcodec", "libx264", "-preset", "fastlibx264", "-preset", "fast")
	args = append(args, "-vf", "hflip", "-copyts", outputFilename)

	rawOutput, err := exec.Command(command, args...).CombinedOutput()
	if err != nil {
		log.Error("Unable exec ffmpeg cmd", err, string(rawOutput[:]))
		return nil, err
	}

	fileContent, err := os.ReadFile(outputFilename)
	if err != nil {
		log.Error("Unable read output file", err)
		return nil, err
	}

	return bytes.NewBuffer(fileContent), err
}
