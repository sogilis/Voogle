package ffmpeg

import (
	"os/exec"
	"strconv"
	"strings"
)

type resolution struct {
	x uint64
	y uint64
}

func (r resolution) GreaterOrEqualResolution(input resolution) bool {
	return r.x >= input.x && r.y >= input.y
}

func CheckContainsSound(filepath string) (bool, error) {
	// sh -c "ffmpeg -i <filepath> 2>&1 | grep Audio | awk '{print $0}' | tr -d ,"
	rawOutput, err := exec.Command("sh", "-c", "ffmpeg -i "+filepath+" 2>&1 | grep Audio | awk '{print $0}' | tr -d ,").CombinedOutput()
	if err != nil {
		return false, err
	}
	haveSound := len(rawOutput) != 0
	return haveSound, err
}

// Extract resolution of the video
func ExtractResolution(filepath string) (resolution, error) {
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
