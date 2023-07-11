package processor

import (
	"context"
	"io"
	"sync"

	"downloader/internal/models"

	"github.com/rs/zerolog"
)

// Downloader Скачивает файлы
type Downloader interface {
	Download(ctx context.Context, task models.Task) (io.ReadCloser, error)
}

type TaskSource interface {
	ExtractTasks(ctx context.Context, url string) (<-chan models.Task, error)
}

// Storage Записывает файлы
type Storage interface {
	WriteFile(ctx context.Context, task models.WriteTask) error
}

// Monitor Врапер для отслеживания прогресса загрузки файла
type Monitor interface {
	Wrap(ctx context.Context, task models.Task, input io.ReadCloser) io.ReadCloser
}

func NewProcessor(
	threads uint,
	downloader Downloader,
	storage Storage,
	taskSource TaskSource,
	monitor Monitor,
) *Processor {
	return &Processor{
		threads:    threads,
		downloader: downloader,
		storage:    storage,
		taskSource: taskSource,
		monitor:    monitor,
	}
}

type Processor struct {
	threads uint

	taskSource TaskSource
	downloader Downloader
	monitor    Monitor
	storage    Storage
}

func (p *Processor) DownloadDirectory(ctx context.Context, url string) error {
	tasks, err := p.taskSource.ExtractTasks(ctx, url)
	if err != nil {
		zerolog.Ctx(ctx).
			Err(err).
			Str("url", url).
			Msg("failed to extract tasks")
		return err
	}

	p.processTasks(ctx, p.threads, tasks)
	return nil
}

func (p *Processor) processTasks(ctx context.Context, count uint, tasks <-chan models.Task) {
	wg := &sync.WaitGroup{}
	wg.Add(int(count))

	for ; count != 0; count-- {
		go p.worker(ctx, wg, tasks)
	}

	wg.Wait()
}

func (p *Processor) worker(ctx context.Context, wg *sync.WaitGroup, tasks <-chan models.Task) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-tasks:
			if !ok {
				return
			}

			err := p.processTask(ctx, task)
			if err != nil {
				zerolog.Ctx(ctx).
					Err(err).
					Str("name", task.Name).
					Str("path", task.Path).
					Str("url", task.URL).
					Msg("failed to process task")
			}
		}
	}
}

func (p *Processor) processTask(ctx context.Context, task models.Task) error {
	file, err := p.downloader.Download(ctx, task)
	if err != nil {
		return err
	}
	defer file.Close()

	// TODO: wrap wireTask, not task
	file = p.monitor.Wrap(ctx, task, file)

	writeTask := models.NewWriteTask(task.Dir, task.Name, file, task.MD5)

	return p.storage.WriteFile(ctx, writeTask)
}
