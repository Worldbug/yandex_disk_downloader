package disk_storage

import (
	"context"
	"crypto/md5"
	"encoding/hex"
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

	if ok, err := ds.md5Verify(task); err == nil || ok {
		return nil
	}

	os.MkdirAll(ds.getPath(task), os.ModePerm)
	out, err := os.Create(path.Join(ds.getPath(task), task.Name))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, task.Source)
	if err != nil {
		return err
	}

	if ok, err := ds.md5Verify(task); err != nil && !ok {
		return errors.Join(err, errors.New("Failed to verify task"))
	}

	return nil
}

func (ds *DiskStorage) getPath(task models.WriteTask) string {
	return ds.root + "/" + task.Dir
}

func (ds *DiskStorage) md5Verify(task models.WriteTask) (bool, error) {
	file, err := os.Open(path.Join(ds.getPath(task), task.Name))
	if err != nil {
		return false, err
	}
	defer file.Close()

	hash := md5.New()

	_, err = io.Copy(hash, file)
	if err != nil {
		return false, err
	}

	sum := hex.EncodeToString(hash.Sum(nil))

	return sum == task.MD5, nil
}
