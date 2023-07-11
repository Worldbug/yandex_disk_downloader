package models

import (
	"net/http"
	"strconv"
)

func NewTask(
	Path string,
	Name string,
	URL string,
	Dir string,
	MD5 string,
) Task {
	return Task{
		Path: Path,
		Name: Name,
		URL:  URL,
		Dir:  Dir,
		MD5:  MD5,
	}
}

type Task struct {
	Path string
	Name string
	URL  string
	Dir  string
	MD5  string
}

func (t *Task) GetSize() int64 {
	resp, err := http.Get(t.URL)
	if err != nil {
		return 0
	}

	size, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	return size
}
