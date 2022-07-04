package ffmpeg

import (
	"bytes"
	"context"
	"io"
	"os/exec"
)

func TransformFlip(ctx context.Context, videoPart []byte) ([]byte, error) {
	return transformHLSPart(ctx, videoPart, []string{"-vf", "vflip"})
}

func TransformGrayscale(ctx context.Context, videoPart []byte) ([]byte, error) {
	return transformHLSPart(ctx, videoPart, []string{"-vf", "hue=s=0"})
}

func transformHLSPart(ctx context.Context, videoPart []byte, ffmepgTransformations []string) ([]byte, error) {
	// Create command
	command := "ffmpeg"
	args := []string{"-i", "pipe:0"}
	args = append(args, "-f", "mpegts", "-muxdelay", "0", "-map", "0:0", "-map", "0:1", "-acodec", "copy")
	args = append(args, "-vcodec", "libx264", "-preset", "fastlibx264", "-preset", "fast", "-copyts")
	args = append(args, ffmepgTransformations...)
	args = append(args, "pipe:1")
	cmd := exec.CommandContext(ctx, command, args...)

	// Fill stdin
	cmd.Stdin = bytes.NewBuffer(videoPart)

	// Connect to stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	// Execute command
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	// Read content of stdout
	grayVideoPart, err := io.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	// Wait end of ffmpeg command
	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	return grayVideoPart, nil
}
