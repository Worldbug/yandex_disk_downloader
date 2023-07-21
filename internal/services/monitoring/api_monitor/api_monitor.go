package api_monitor

import (
	"context"
	"io"

	"downloader/internal/models"
)

func NewStreamMonitor() *StreamMonitor {
	return &StreamMonitor{}
}

type StreamMonitor struct{}

func (sm *StreamMonitor) Wrap(
	ctx context.Context,
	task models.Task,
	input io.ReadCloser,
) io.ReadCloser {
	// TODO:
	return input
}
