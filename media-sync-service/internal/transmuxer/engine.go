package transmuxer

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Transmuxer struct{}

func NewTransmuxer() *Transmuxer {
	return &Transmuxer{}
}

func (t *Transmuxer) Transmux(ctx context.Context, inputPath, outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, os.FileMode(0755)); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}
	masterPlaylistPath := filepath.Join(outputDir, "master.m3u8")
	segmentPattern := filepath.Join(outputDir, "segment_%03d.ts")
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-y",
		"-i", inputPath,
		"-c", "copy",
		"-hls_time", "6",
		"-hls_playlist_type", "vod",
		"-hls_flags", "independent_segments",
		"-hls_segment_filename", segmentPattern,
		masterPlaylistPath,
	)
	out,err:= cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg command failed: %v, output: %s", err, string(out))
	}
	return masterPlaylistPath, nil
}
