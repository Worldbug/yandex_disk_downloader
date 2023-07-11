package cli_monitor

import (
	"context"
	"io"

	"downloader/internal/models"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

func NewCliMonitor() *CliMonitor {
	return &CliMonitor{
		progressBar: mpb.New(),
	}
}

type CliMonitor struct {
	progressBar *mpb.Progress
}

func (cm *CliMonitor) Wrap(
	ctx context.Context,
	task models.Task,
	input io.ReadCloser,
) io.ReadCloser {
	size := task.GetSize()
	bar := cm.addBar(task.Name, size)
	lr := io.LimitReader(input, size)
	return bar.ProxyReader(lr)
}

func (cm *CliMonitor) addBar(name string, total int64) *mpb.Bar {
	return cm.progressBar.AddBar(total,
		mpb.BarRemoveOnComplete(),
		mpb.PrependDecorators(
			decor.Name(name),
			decor.Percentage(decor.WCSyncSpace),
			decor.CountersKibiByte(" % .2f / % .2f "),
		),
		mpb.AppendDecorators(
			decor.EwmaSpeed(decor.UnitKiB, "% .1f ", 60),
			decor.OnComplete(
				decor.EwmaETA(decor.ET_STYLE_GO, 60, decor.WCSyncWidth), "done",
			),
		),
	)
}
