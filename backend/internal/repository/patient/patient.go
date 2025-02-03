package repository

import (
	"context"
	"dicomviewer/internal/models"
)

type PatientRepository interface {
	Insert(ctx context.Context, data models.DicomFileMetadata) error
	List(ctx context.Context, page models.PageQuery) (interface{}, error)
}
