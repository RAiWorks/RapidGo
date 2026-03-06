package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// FileStore stores one JSON file per session in the configured directory.
type FileStore struct {
	Path string
}

func (s *FileStore) filepath(id string) string {
	return filepath.Join(s.Path, id+".json")
}

func (s *FileStore) Read(id string) (map[string]interface{}, error) {
	raw, err := os.ReadFile(s.filepath(id))
	if err != nil {
		return make(map[string]interface{}), nil
	}
	var data map[string]interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *FileStore) Write(id string, data map[string]interface{}, lifetime time.Duration) error {
	if err := os.MkdirAll(s.Path, 0700); err != nil {
		return err
	}
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(s.filepath(id), raw, 0600)
}

func (s *FileStore) Destroy(id string) error {
	return os.Remove(s.filepath(id))
}

func (s *FileStore) GC(maxLifetime time.Duration) error {
	entries, err := os.ReadDir(s.Path)
	if err != nil {
		return err
	}
	cutoff := time.Now().Add(-maxLifetime)
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			os.Remove(filepath.Join(s.Path, e.Name()))
		}
	}
	return nil
}
