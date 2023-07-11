package api_monitor

import (
	"context"
	"io"

	"downloader/internal/models"

	"go.uber.org/ratelimit"
)

type StreamMonitor struct {
	taskMap map[string]*discoverTask
}

func (sm *StreamMonitor) Wrap(
	ctx context.Context,
	task models.Task,
	input io.ReadCloser,
) io.ReadCloser {
	dt := NewDiscoverTask()
	dt.onChange = sm.onChange

	return nil
}

func (sm *StreamMonitor) addTask(ctx context.Context, task models.Task) {}

func (sm *StreamMonitor) onChange() {
	rl := ratelimit.New(1)
}

func NewDiscoverTask() *discoverTask {
	return &discoverTask{}
}

type discoverTask struct {
	FileName string
	Size     int64
	Loaded   int64
	onChange func()
}
