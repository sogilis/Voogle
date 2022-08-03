package ffmpeg

import (
	"context"
	"io"
	"os/exec"
)

func TransformFlip(ctx context.Context, videoPart io.Reader, transformedVideoPart io.Writer) error {
	return transformHLSPart(ctx, videoPart, transformedVideoPart, []string{"-vf", "vflip"})
}

func TransformGrayscale(ctx context.Context, videoPart io.Reader, transformedVideoPart io.Writer) error {
	return transformHLSPart(ctx, videoPart, transformedVideoPart, []string{"-vf", "hue=s=0"})
}

func transformHLSPart(ctx context.Context, videoPart io.Reader, transformedVideoPart io.Writer, ffmepgTransformations []string) error {
	// Create command
	command := "ffmpeg"
	args := []string{"-i", "pipe:0"}
	args = append(args, "-f", "mpegts", "-muxdelay", "0", "-map", "0:0", "-map", "0:1", "-acodec", "copy")
	args = append(args, "-vcodec", "libx264", "-preset", "fastlibx264", "-preset", "superfast", "-copyts")
	args = append(args, ffmepgTransformations...)
	args = append(args, "pipe:1")
	cmd := exec.CommandContext(ctx, command, args...)

	// Fill stdin
	cmd.Stdin = videoPart

	cmd.Stdout = transformedVideoPart
	// Execute command
	err := cmd.Start()
	if err != nil {
		return err
	}

	// Wait end of ffmpeg command
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}
