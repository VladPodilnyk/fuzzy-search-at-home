package main

import (
	"dicomviewer/internal/app"
	"dicomviewer/internal/service"
)

func main() {
	imageService := service.NewImageService()
	viewerService := service.NewViewerService(nil)

	application := app.New(imageService, viewerService)
	application.Run()
}
