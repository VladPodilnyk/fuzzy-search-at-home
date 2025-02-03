package app

import (
	"context"
	"dicomviewer/internal/service"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

type AppData struct {
	imageService  service.ImageService
	viewerService service.ViewerService
}

func New(imageService service.ImageService, viewerService service.ViewerService) *AppData {
	return &AppData{imageService: imageService, viewerService: viewerService}
}

func (app *AppData) Run() {
	e := app.getRoutes()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Launch server
	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait here for an interrupt signal
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func (app *AppData) getRoutes() *echo.Echo {
	mux := echo.New()

	mux.GET("/list", app.viewerService.GetRecords)
	mux.GET("/download", app.imageService.DownloadFile)
	mux.GET("/preview", app.imageService.GetFilePreview)
	mux.POST("/upload", app.imageService.UploadFile)
	return mux
}
