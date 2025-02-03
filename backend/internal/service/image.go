package service

import (
	"errors"

	"github.com/labstack/echo/v4"
)

type ImageService interface {
	UploadFile(ctx echo.Context) error
	DownloadFile(ctx echo.Context) error
	GetFilePreview(ctx echo.Context) error
}

type ImageServiceImpl struct{}

func NewImageService() ImageServiceImpl {
	return ImageServiceImpl{}
}

func (impl ImageServiceImpl) UploadFile(ctx echo.Context) error {
	return errors.New("not implemented")
}
func (impl ImageServiceImpl) DownloadFile(ctx echo.Context) error {
	return errors.New("not implemented")
}
func (impl ImageServiceImpl) GetFilePreview(ctx echo.Context) error {
	return errors.New("not implemented")
}
