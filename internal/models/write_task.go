package models

import "io"

func NewWriteTask(
	dir string,
	name string,
	source io.ReadCloser,
	md5 string,
) WriteTask {
	return WriteTask{
		Dir:    dir,
		Name:   name,
		Source: source,
		MD5:    md5,
	}
}

type WriteTask struct {
	Dir    string
	Name   string
	Source io.ReadCloser
	MD5    string
}

func (wt *WriteTask) Done() {
	wt.Source.Close()
}
