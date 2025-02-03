package service

import (
	repository "dicomviewer/internal/repository/files"
	"errors"

	"github.com/labstack/echo/v4"
)

type ViewerService interface {
	GetRecords(ctx echo.Context) error
}

type ViewerServiceImpl struct {
	repo repository.FilesRepository
}

func NewViewerService(repo repository.FilesRepository) ViewerServiceImpl {
	return ViewerServiceImpl{repo: repo}
}

func (s ViewerServiceImpl) GetRecords(ctx echo.Context) error {
	return errors.New("not implemented")
}
