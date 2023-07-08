package disk_storage

import (
	"context"
	"crypto/md5"
	"errors"
	"io"
	"os"
	"path"

	"downloader/internal/models"
)

func NewDiskStorage(root string) *DiskStorage {
	return &DiskStorage{
		root: root,
	}
}

type DiskStorage struct {
	root string
}

func (ds *DiskStorage) WriteFile(
	ctx context.Context,
	task models.WriteTask,
) error {
	defer task.Done()

	if ds.md5Verify(task) {
		return nil
	}

	os.MkdirAll(ds.getPath(task), os.ModePerm)
	out, err := os.Create(path.Join(task.Dir, task.Name))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, task.Source)
	if err != nil {
		return err
	}

	if !ds.md5Verify(task) {
		return errors.New("Failed to verify task")
	}

	return nil
}

func (ds *DiskStorage) getPath(task models.WriteTask) string {
	return ds.root + "/" + task.Dir
}

func (ds *DiskStorage) md5Verify(task models.WriteTask) bool {
	file, err := os.Open(ds.getPath(task))
	if err != nil {
		return false
	}
	defer file.Close()

	hash := md5.New()

	_, err = io.Copy(hash, file)
	if err != nil {
		return false
	}

	return string(hash.Sum(nil)) == task.MD5
}
