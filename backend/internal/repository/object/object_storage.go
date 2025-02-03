package repository

import "context"

type ObjectStorage interface {
	Save(ctx context.Context, file interface{}) (string, error)
	Get(ctx context.Context, filepath string) (interface{}, error)
}
