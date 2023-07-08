package yandex_disk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"downloader/internal/models"

	"github.com/rs/zerolog"
)

const (
	base_url = "https://cloud-api.yandex.net/v1/disk/public/resources"
	rootPath = ""
)

func newYandexDiskClient() *YandexDiskClient {
	return &YandexDiskClient{}
}

type YandexDiskClient struct{}

func (ydc *YandexDiskClient) ExtractTasks(
	ctx context.Context,
	url string,
) (<-chan models.Task, error) {
	tasks := make(chan models.Task)

	go func() {
		defer close(tasks)

		select {
		case <-ctx.Done():
			return
		default:
			err := GetFileTree(ctx, url, tasks)
			if err != nil {
				zerolog.Ctx(ctx).
					Err(err).
					Str("url", url).
					Msg("Failed to get file tree")
			}
		}
	}()

	return tasks, nil
}

func GetFileTree(
	ctx context.Context,
	url string,
	tasks chan<- models.Task,
) error {
	return introspectPath(url, rootPath, rootPath, tasks)
}

func introspectPath(url, path, dir string, tasks chan<- models.Task) error {
	resp, err := describe(url, path)
	if err != nil {
		return err
	}

	for _, item := range resp.Embedded.Items {
		if item.Type == "file" {
			task := models.NewTask(
				item.Path,
				item.Name,
				item.File,
				dir,
				item.MD5,
			)

			tasks <- task
		}

		if item.Type == "dir" {
			introspectPath(
				url,
				item.Path,
				fmt.Sprint(dir, item.Name, "/"),
				tasks,
			)
		}
	}

	return nil
}

func describe(key, path string) (*YandexResp, error) {
	base, _ := url.Parse(base_url)
	params := url.Values{}
	params.Add("public_key", key)
	params.Add("path", path)
	base.RawQuery = params.Encode()

	yResp := &YandexResp{}
	resp, err := http.Get(base.String())
	if err != nil {
		return yResp, err
	}

	err = json.NewDecoder(resp.Body).Decode(yResp)
	if err != nil {
		return yResp, err
	}

	return yResp, nil
}
