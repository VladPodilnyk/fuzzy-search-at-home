package repository

import (
	"context"
	"dicomviewer/internal/models"
	"fmt"

	"github.com/google/uuid"
)

type FilesRepository interface {
	Insert(ctx context.Context, data models.FileData) error
	Get(ctx context.Context, id uint) (models.FileData, error)
	List(ctx context.Context, page models.PageQuery) ([]models.FileDataRow, error)
	Count(ctx context.Context) (uint, error)
}

type InMemFilesRepository struct {
	store map[uint]models.FileData
}

func NewInMemFilesRepository() *InMemFilesRepository {
	return &InMemFilesRepository{store: make(map[uint]models.FileData)}
}

func (impl *InMemFilesRepository) Insert(ctx context.Context, data models.FileData) error {
	newId := uint(uuid.New().ID())
	impl.store[newId] = data
	return nil
}

func (impl *InMemFilesRepository) Get(ctx context.Context, id uint) (models.FileData, error) {
	value, isOk := impl.store[id]
	if !isOk {
		return models.FileData{}, fmt.Errorf("cannot find file attachment with id %d", id)
	}
	return value, nil
}

func (impl *InMemFilesRepository) List(ctx context.Context, page models.PageQuery) ([]models.FileDataRow, error) {
	res := make([]models.FileDataRow, 0, page.Limit)
	offset := uint(0)
	for key, value := range impl.store {
		if offset < page.Offset {
			offset += 1
			continue
		}

		if len(res) == int(page.Limit) {
			break
		}

		res = append(res, models.FileDataRow{
			Id:                key,
			OriginalFile:      value.OriginalFile,
			PreviewFile:       value.PreviewFile,
			SeriesDescription: value.SeriesDescription,
			UserId:            value.UserId,
		})
	}
	return res, nil
}

func (impl *InMemFilesRepository) Count(ctx context.Context) (uint, error) {
	return uint(len(impl.store)), nil
}
