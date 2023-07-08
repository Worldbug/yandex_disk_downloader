package processor

import (
	"context"
	"io"
	"sync"

	"downloader/internal/models"
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

// Status Интерфейс для вывод статуса обработки задач
type Status interface{}

func NewProcessor(
	rootURL string,
	threads uint,
	downloader Downloader,
	storage Storage,
	status Status,
) *Processor {
	return &Processor{
		rootURL:    rootURL,
		threads:    threads,
		downloader: downloader,
		storage:    storage,
		status:     status,
	}
}

type Processor struct {
	rootURL string
	threads uint

	taskSource TaskSource
	downloader Downloader
	storage    Storage
	status     Status

	cancel context.CancelFunc
}

func (p *Processor) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	p.cancel = cancel

	tasks, err := p.taskSource.ExtractTasks(ctx, p.rootURL)
	if err != nil {
		// TODO: err log
		return
	}

	p.processTasks(ctx, p.threads, tasks)
}

func (p *Processor) Stop() {
	if p.cancel != nil {
		p.cancel()
	}
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
				// TODO: error log
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

	writeTask := models.NewWriteTask(task.Path, task.Name, file, task.MD5)

	return p.storage.WriteFile(ctx, writeTask)
}
