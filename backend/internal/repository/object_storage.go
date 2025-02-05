package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	ORIGINAL = "original"
	PREVIEW  = "preview"
)

type ObjectStorage interface {
	Save(ctx context.Context, path string, data []byte) error
	Get(ctx context.Context, path string) ([]byte, error)
}

type InMemoryObjStorage struct {
	storage map[string][]byte
}

func NewInMemoryObjectStorage() *InMemoryObjStorage {
	return &InMemoryObjStorage{storage: make(map[string][]byte)}
}

func MakeFilePath(prefix, filename string) string {
	now := time.Now().Format("2006-01-02")
	return fmt.Sprintf("%s/%s/%s_%s", prefix, now, uuid.New().String(), filename)
}

func (s *InMemoryObjStorage) Save(ctx context.Context, path string, data []byte) error {
	s.storage[path] = data
	return nil
}

func (s *InMemoryObjStorage) Get(ctx context.Context, path string) ([]byte, error) {
	value, ok := s.storage[path]
	if !ok {
		return nil, fmt.Errorf("cannot find a file by the given path: %s", path)
	}
	return value, nil
}
