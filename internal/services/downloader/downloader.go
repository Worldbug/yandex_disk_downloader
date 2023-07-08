package downloader

import (
	"context"
	"io"
	"net/http"

	"downloader/internal/models"
)

func NewHTTPDownloader() *HTTPDownloader {
	return &HTTPDownloader{}
}

type HTTPDownloader struct{}

func (d *HTTPDownloader) Download(ctx context.Context, task models.Task) (io.ReadCloser, error) {
	resp, err := http.Get(task.URL)
	return resp.Body, err
}
