package main

import (
	_ "dicomviewer/cmd/docs"
	"dicomviewer/internal/app"
	"dicomviewer/internal/repository"
	"dicomviewer/internal/service"
)

func main() {
	// Init repository layer
	objStorage := repository.NewInMemoryObjectStorage()
	filesStorage := repository.NewInMemFilesRepository()
	patientStorage := repository.NewInMemPatientRepository()

	// Init service layer
	imageService := service.NewImageService(objStorage, patientStorage, filesStorage)

	// Launch
	application := app.New(imageService)
	application.Run()
}
