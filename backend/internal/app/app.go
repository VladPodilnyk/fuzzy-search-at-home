package app

import (
	"context"
	"dicomviewer/internal/service"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type AppData struct {
	imageService service.ImageService
}

func New(imageService service.ImageService) *AppData {
	return &AppData{imageService: imageService}
}

func (app *AppData) Run() {
	e := app.getRoutes()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Launch server
	go func() {
		if err := e.Start("localhost:1323"); err != nil && err != http.ErrServerClosed {
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

// @Title		My Echo API
// @Version		1.0
// @Description	This is a sample API using Swaggo with Echo
// @BasePath	/
func (app *AppData) getRoutes() *echo.Echo {
	mux := echo.New()

	mux.POST("/list", app.imageService.ListFiles)
	mux.POST("/download", app.imageService.DownloadFile)
	mux.POST("/preview", app.imageService.GetFilePreview)
	mux.POST("/upload", app.imageService.UploadFile)

	mux.GET("/swagger/*", echoSwagger.WrapHandler)
	return mux
}
