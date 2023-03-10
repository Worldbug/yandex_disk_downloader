// Simple script for download large folders form yandex disk
// Usage: go run . [URL] [THREAD Count]
// Example: go run . https://disk.yandex.ru/d/foler_url 3
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

const base_url = "https://cloud-api.yandex.net/v1/disk/public/resources"

func main() {
	if len(os.Args) < 3 {
		log.Println("Not enoughs args ")
		return
	}

	u, err := url.Parse(os.Args[1])
	if err != nil {
		log.Println("Wrong url format: ", err)
		return
	}

	threads, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Println("Can`t parse thread count arg: ", err)
		return
	}

	progressBar = mpb.New()
	newYandexDownloader(u.String(), threads).start()
}

func (yd *YandexDownloader) ride(url, path, dir string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	resp := describe(url, path)
	for _, item := range resp.Embedded.Items {
		if item.Type == "file" {
			yd.addTask(Task{
				path: item.Path,
				name: item.Name,
				url:  item.File,
				dir:  dir,
			})
		}

		if item.Type == "dir" {
			yd.ride(
				url,
				item.Path,
				fmt.Sprint(dir, item.Name, "/"),
				wg,
			)
		}
	}
}

func describe(key, path string) *YandexResp {
	base, _ := url.Parse(base_url)
	params := url.Values{}
	params.Add("public_key", key)
	params.Add("path", path)
	base.RawQuery = params.Encode()

	resp, err := http.Get(base.String())
	if err != nil {

	}

	yResp := &YandexResp{}

	err = json.NewDecoder(resp.Body).Decode(yResp)
	if err != nil {

	}

	return yResp
}

type YandexResp struct {
	Embedded struct {
		Sort      string `json:"sort"`
		PublicKey string `json:"public_key"`
		Items     []Item `json:"items"`
		Limit     int    `json:"limit"`
		Offset    int    `json:"offset"`
		Path      string `json:"path"`
		Total     int    `json:"total"`
	} `json:"_embedded"`
}

type Item struct {
	Path string `json:"path"`
	Type string `json:"type"`
	Name string `json:"name"`
	File string `json:"file,omitempty"`
}

var progressBar *mpb.Progress

func newBar(name string, total int64) *mpb.Bar {
	return progressBar.AddBar(total,
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

func newYandexDownloader(url string, workers int) *YandexDownloader {
	return &YandexDownloader{
		url:     url,
		tasks:   make(chan Task),
		workers: workers,
	}
}

type YandexDownloader struct {
	url     string
	tasks   chan Task
	workers int
}

func (yd *YandexDownloader) start() {
	for ; yd.workers > 0; yd.workers-- {
		go yd.worker()
	}

	wg := &sync.WaitGroup{}
	yd.ride(yd.url, "", "", wg)
	wg.Wait()
}

func (yd *YandexDownloader) worker() {
	for task := range yd.tasks {
		loader(task)
	}
}

func (yd *YandexDownloader) addTask(t Task) {
	yd.tasks <- t
}

type Task struct {
	path string
	name string
	url  string
	dir  string
}

func loader(task Task) error {
	os.MkdirAll("./"+task.dir, os.ModePerm)
	out, err := os.Create(path.Join(task.dir, task.name))
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(task.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	size, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)

	r := io.LimitReader(resp.Body, size)
	bar := newBar(task.name, size)

	if _, err = io.Copy(out, bar.ProxyReader(r)); err != nil {
		return err
	}

	return err
}
