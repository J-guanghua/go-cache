package store

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 实现本地文件缓存接口
// Implement local file cache interface
type fileStore struct {
	level     int    // Optimize the file path level
	directory string // cache directory
	fileMode  os.FileMode
}

type FileOption func(store *fileStore)

func FileMode(mode os.FileMode) FileOption {
	return func(store *fileStore) {
		store.fileMode = mode
	}
}

func Level(levle int) FileOption {
	return func(store *fileStore) {
		store.level = levle
	}
}

func Directory(cachePath string) FileOption {
	return func(store *fileStore) {
		store.directory = cachePath
	}
}

func NewFile(opts ...FileOption) Store {
	fileCache := &fileStore{
		directory: "./tmp",
		fileMode:  755,
		level:     2,
	}
	for _, o := range opts {
		o(fileCache)
	}
	_, err := os.Stat(fileCache.directory)
	if os.IsNotExist(err) {
		err = os.Mkdir(fileCache.directory, fileCache.fileMode)
		if err != nil {
			panic(err)
		}
	}
	return fileCache
}

func (file *fileStore) Get(ctx context.Context, name string) ([]byte, error) {
	name = file.buildFile(ctx, name)
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fi.ModTime().Second() > time.Now().Second() {
		return ioutil.ReadAll(f)
	}
	return nil, ErrNotFound
}

func (file *fileStore) Set(ctx context.Context, name string, v []byte, expiration time.Duration) error {
	name = file.buildFile(ctx, name)
	err := os.WriteFile(name, v, file.fileMode)
	if err != nil {
		return err
	}
	return os.Chtimes(name, time.Now(), time.Now().Add(expiration))
}

func (file *fileStore) Del(ctx context.Context, name string) error {
	name = file.buildFile(ctx, name)
	return os.Remove(name)
}

func (file *fileStore) buildFile(_ context.Context, name string) string {
	dir := filepath.Join(file.directory, name[len(name)-file.level:])
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		_ = os.Mkdir(dir, file.fileMode)
	}
	return filepath.Join(dir, name)
}

// Matching the filename prefix clears the file
func (file *fileStore) Flush(_ context.Context, prefix string) error {
	return treeFile(file.directory, func(fileName string) error {
		_, name := filepath.Split(fileName)
		if strings.HasPrefix(name, prefix) {
			return os.Remove(fileName)
		}
		return nil
	})
}

func (file *fileStore) Gc(_ context.Context) error {
	return nil
}

// Iterate through the cache file recursively
func treeFile(cachePath string, rm func(fileName string) error) error {
	files, err := ioutil.ReadDir(cachePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			err := treeFile(filepath.Join(cachePath, file.Name()), rm)
			if err != nil {
				return err
			}
		} else {
			err := rm(filepath.Join(cachePath, file.Name()))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
