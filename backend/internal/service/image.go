package service

import (
	"bytes"
	"context"
	"dicomviewer/internal/dicomutils"
	"dicomviewer/internal/models"
	"dicomviewer/internal/repository"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type ImageService interface {
	UploadFile(ctx echo.Context) error
	DownloadFile(ctx echo.Context) error
	GetFilePreview(ctx echo.Context) error
	ListFiles(ctx echo.Context) error
}

type ImageServiceImpl struct {
	objectStorage repository.ObjectStorage
	patientRepo   repository.PatientRepository
	filesRepo     repository.FilesRepository
}

func NewImageService(
	objectStorage repository.ObjectStorage,
	patientRepo repository.PatientRepository,
	filesRepo repository.FilesRepository,
) *ImageServiceImpl {
	return &ImageServiceImpl{
		objectStorage: objectStorage,
		patientRepo:   patientRepo,
		filesRepo:     filesRepo,
	}
}

func (impl *ImageServiceImpl) UploadFile(ctx echo.Context) error {
	context := ctx.Request().Context()

	header, err := getFileHeader(ctx)
	if err != nil {
		return logError(err)
	}

	file, err := header.Open()
	if err != nil {
		return logError(err)
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return logError(err)
	}

	fileData, err := impl.extractData(bytes.NewReader(fileBytes), header.Size)
	if err != nil {
		return logError(err)
	}

	userId, err := impl.patientRepo.Insert(context, fileData.Meta.PatientName, fileData.Meta.BirthDate)
	if err != nil {
		return logError(err)
	}

	pathInfo, err := impl.storeFiles(context, header.Filename, fileBytes, fileData.Image)
	if err != nil {
		return logError(err)
	}

	err = impl.storeFileData(context, userId, pathInfo, fileData.Meta.SeriesDescription)
	if err != nil {
		return logError(err)
	}

	return ctx.JSON(http.StatusOK, nil)
}

func (impl *ImageServiceImpl) DownloadFile(ctx echo.Context) error {
	var fileAttachmentId models.FileAttachmentId

	if err := ctx.Bind(&fileAttachmentId); err != nil {
		return logError(err)
	}

	context := ctx.Request().Context()
	row, err := impl.filesRepo.Get(context, fileAttachmentId.Id)
	if err != nil {
		return logError(err)
	}

	bytes, err := impl.objectStorage.Get(context, row.OriginalFile)
	if err != nil {
		return logError(err)
	}

	return ctx.Blob(http.StatusOK, "application/dicom", bytes)
}

// TODO: make a generic function for file download
func (impl *ImageServiceImpl) GetFilePreview(ctx echo.Context) error {
	var fileAttachmentId models.FileAttachmentId

	if err := ctx.Bind(&fileAttachmentId); err != nil {
		return logError(err)
	}

	context := ctx.Request().Context()
	row, err := impl.filesRepo.Get(context, fileAttachmentId.Id)
	if err != nil {
		return logError(err)
	}

	bytes, err := impl.objectStorage.Get(context, row.PreviewFile)
	if err != nil {
		return logError(err)
	}

	return ctx.Blob(http.StatusOK, "application/png", bytes)
}

//	@Summary		List files
//	@Description	get all file attachments
//	@Accept			json
//	@Produce		json
//	@Param			page	body	models.PageQuery	true	"Page start and size"
//	@Success		200		{array}	models.FileDataRow	"List of file attachments"
//	@Router			/list [post]
func (impl *ImageServiceImpl) ListFiles(ctx echo.Context) error {
	context := ctx.Request().Context()

	var page models.PageQuery
	if err := ctx.Bind(&page); err != nil {
		return logError(err)
	}

	// TODO: rework list stuff
	data, err := impl.filesRepo.List(context, page)
	if err != nil {
		return logError(err)
	}

	return ctx.JSON(http.StatusOK, data)
}

func (impl *ImageServiceImpl) storeFileData(ctx context.Context, userId uint, pathInfo *models.FilePath, series string) error {
	row := models.FileData{
		UserId:            userId,
		OriginalFile:      pathInfo.Original,
		PreviewFile:       pathInfo.Preview,
		SeriesDescription: series,
	}
	return impl.filesRepo.Insert(ctx, row)
}

func (impl *ImageServiceImpl) storeFiles(ctx context.Context, filename string, originalFile []byte, preview image.Image) (*models.FilePath, error) {
	orignalFilePath, err := impl.uploadFile(ctx, filename, originalFile, true)
	if err != nil {
		return nil, err
	}

	previewImagePath, err := impl.uploadPreviewImage(ctx, filename, preview)
	if err != nil {
		return nil, err
	}
	return &models.FilePath{Original: orignalFilePath, Preview: previewImagePath}, nil
}

func (impl *ImageServiceImpl) uploadPreviewImage(ctx context.Context, filename string, img image.Image) (string, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return "", err
	}

	nameWithoutExtension := strings.Split(filename, ".")[0]
	return impl.uploadFile(ctx, nameWithoutExtension, buf.Bytes(), false)
}

func (impl *ImageServiceImpl) uploadFile(ctx context.Context, filename string, data []byte, isOriginal bool) (string, error) {
	var prefix string
	if isOriginal {
		prefix = repository.ORIGINAL
	} else {
		prefix = repository.PREVIEW
	}

	path := repository.MakeFilePath(prefix, filename)
	err := impl.objectStorage.Save(ctx, path, data)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (impl *ImageServiceImpl) extractData(in io.Reader, size int64) (*models.DicomFileInfo, error) {
	dataset, err := dicomutils.ReadDicomDataset(in, size)
	if err != nil {
		return nil, err
	}

	metadata, err := dicomutils.GetMetadata(dataset)
	if err != nil {
		return nil, err
	}

	img, err := dicomutils.GetImage(dataset)
	if err != nil {
		return nil, err
	}

	return &models.DicomFileInfo{Meta: *metadata, Image: img}, err
}

func getFileHeader(ctx echo.Context) (*multipart.FileHeader, error) {
	file, err := ctx.FormFile("file")
	if err != nil {
		return nil, err
	}
	return file, err
}
