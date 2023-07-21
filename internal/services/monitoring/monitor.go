package monitoring

import (
	"context"
	"io"
)

/*
	- Скорость загрузки
	- Скаченные файлы
	- Файлы в процессе скачивания
		- Имя
		- Скорость
		- Процент
*/

type Monitor interface {
	WrapWithMonitoring(ctx context.Context, task any) any
	GetDownloadedList(ctx context.Context) []string
	GetDownloadFileStream(ctx context.Context)
}

type FileStatus struct {
	Name          string
	MbPerSecounds int
	stream        io.ReadCloser
}
