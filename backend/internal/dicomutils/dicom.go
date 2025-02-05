package dicomutils

import (
	"dicomviewer/internal/models"
	"fmt"
	"image"
	"image/color"
	"io"
	"math"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

func ReadDicomDataset(in io.Reader, fileSize int64) (*dicom.Dataset, error) {
	dataset, err := dicom.Parse(in, fileSize, nil)
	if err != nil {
		return nil, fmt.Errorf("error parsing DICOM file: %w", err)
	}
	return &dataset, nil
}

func GetMetadata(ds *dicom.Dataset) (*models.DicomFileMetadata, error) {
	element, err := ds.FindElementByTag(tag.PatientName)
	if err != nil {
		return nil, fmt.Errorf("couldn't extract patient name: %w", err)
	}
	patientName := dicom.MustGetStrings(element.Value)[0]

	element, err = ds.FindElementByTag(tag.PatientBirthDate)
	if err != nil {
		return nil, fmt.Errorf("couldn't get patient birth date: %w", err)
	}
	birthDate := dicom.MustGetStrings(element.Value)[0]

	element, err = ds.FindElementByTag(tag.SeriesDescription)
	if err != nil {
		return nil, fmt.Errorf("couldn't get series description: %w", err)
	}
	seriesDescription := dicom.MustGetStrings(element.Value)[0]

	result := models.DicomFileMetadata{
		PatientName:       patientName,
		BirthDate:         birthDate,
		SeriesDescription: seriesDescription,
	}

	return &result, nil
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

// Not ideal, but very simple solution that
// does more or less good job
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
