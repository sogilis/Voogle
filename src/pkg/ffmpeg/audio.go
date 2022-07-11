package ffmpeg

import (
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func AddEmptyAudioTrack(fileName string) error {
	// ffmpeg -f lavfi -i anullsrc=channel_layout=stereo:sample_rate=44100 -i <filepath> -c:v copy -c:a aac -shortest <filepath>
	tmpPath := "tmp" + filepath.Ext(fileName)
	_, err := exec.Command("ffmpeg", "-y", "-f", "lavfi", "-i", "anullsrc=channel_layout=stereo:sample_rate=44100", "-i", fileName, "-c:v", "copy", "-c:a", "aac", "-shortest", tmpPath).CombinedOutput()
	if err != nil {
		return err
	}
	log.Debug("Empty audio generated")
	return os.Rename(tmpPath, fileName)
}
