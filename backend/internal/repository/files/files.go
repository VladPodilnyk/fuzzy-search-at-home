package repository

import (
	"context"
	"dicomviewer/internal/models"
)

type FilesRepository interface {
	Insert(ctx context.Context) error
	List(ctx context.Context, page models.PageQuery) (interface{}, error)
}
