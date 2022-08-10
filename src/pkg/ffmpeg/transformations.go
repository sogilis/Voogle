package ffmpeg

import (
	"context"
	"io"
	"os/exec"
)

func CreateFlipCommand(ctx context.Context) *exec.Cmd {
	// Create command
	command := "ffmpeg"
	args := []string{"-i", "pipe:0"}
	args = append(args, "-f", "mpegts", "-muxdelay", "0", "-map", "0:0", "-map", "0:1", "-acodec", "copy")
	args = append(args, "-vcodec", "libx264", "-preset", "fastlibx264", "-preset", "superfast", "-copyts")
	args = append(args, "-vf", "vflip")
	args = append(args, "pipe:1")
	return exec.CommandContext(ctx, command, args...)
}

func CreateGrayCommand(ctx context.Context) *exec.Cmd {
	// Create command
	command := "ffmpeg"
	args := []string{"-i", "pipe:0"}
	args = append(args, "-f", "mpegts", "-muxdelay", "0", "-map", "0:0", "-map", "0:1", "-acodec", "copy")
	args = append(args, "-vcodec", "libx264", "-preset", "fastlibx264", "-preset", "superfast", "-copyts")
	args = append(args, "-vf", "hue=s=0")
	args = append(args, "pipe:1")
	return exec.CommandContext(ctx, command, args...)
}

func TransformHLSPart(cmd *exec.Cmd, stdin io.Reader, stdout io.Writer) error {
	cmd.Stdin = stdin
	cmd.Stdout = stdout

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
