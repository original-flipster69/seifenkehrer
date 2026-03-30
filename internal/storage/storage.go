package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

type Storage interface {
	Load(path string) ([]byte, error)
	Save(path string, data []byte) error
}

type FileStorage struct{}

func (FileStorage) Load(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	return data, nil
}

func (FileStorage) Save(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0775); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	tmp, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return fmt.Errorf("writing %s: %w", path, err)
	}

	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return fmt.Errorf("closing temp file: %w", err)
	}

	return os.Rename(tmp.Name(), path)
}
