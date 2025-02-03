package dicom

import (
	"dicomviewer/internal/models"
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

func ReadDicomDataset(filepath string) (*dicom.Dataset, error) {
	dataset, err := dicom.ParseFile(filepath, nil)
	if err != nil {
		return nil, fmt.Errorf("error parsing DICOM file: %w", err)
	}
	return &dataset, nil
}

func GetMetadata(ds *dicom.Dataset) (*models.DicomFileMetadata, error) {
	return nil, nil
}

func GetImage(ds *dicom.Dataset) (image.Image, error) {
	pixelData, err := ds.FindElementByTag(tag.PixelData)
	if err != nil {
		return nil, fmt.Errorf("error finding pixel data: %w", err)
	}
	pixelDataInfo := dicom.MustGetPixelDataInfo(pixelData.Value)
	frame, err := pixelDataInfo.Frames[0].GetImage()
	if err != nil {
		return nil, err
	}
	img := scaleIntensity(frame)
	return img, nil
}

func scaleIntensity(img image.Image) image.Image {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	var min, max uint8 = 255, 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayImg.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
			pixel := grayImg.GrayAt(x, y)
			if pixel.Y < min {
				min = pixel.Y
			}
			if pixel.Y > max {
				max = pixel.Y
			}
		}
	}

	scale := 255.0 / math.Max(1, float64(max-min))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := grayImg.GrayAt(x, y)
			scaled := uint8(math.Round(float64(pixel.Y-min) * scale))
			grayImg.SetGray(x, y, color.Gray{Y: scaled})
		}
	}

	return grayImg
}
