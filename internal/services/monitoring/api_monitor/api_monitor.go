package api_monitor

import (
	"context"
	"io"

	"downloader/internal/models"
)

type StreamMonitor struct{}

func (sm *StreamMonitor) Wrap(
	ctx context.Context,
	task models.Task,
	input io.ReadCloser,
) io.ReadCloser {
	// TODO:
	return input
}
