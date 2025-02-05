package models

import "image"

type DicomFileMetadata struct {
	PatientName       string `json:"patient_name"`
	BirthDate         string `json:"birth_date"`
	SeriesDescription string `json:"series_description"`
}

type UserData struct {
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
}

type DicomFileInfo struct {
	Meta  DicomFileMetadata
	Image image.Image
}

type FilePath struct {
	Original string
	Preview  string
}

type FileData struct {
	OriginalFile      string `json:"original_file"`
	PreviewFile       string `json:"preview_file"`
	SeriesDescription string `json:"series_description"`
	UserId            uint   `json:"user_id"`
}

type FileDataRow struct {
	Id                uint   `json:"id"`
	OriginalFile      string `json:"original_file"`
	PreviewFile       string `json:"preview_file"`
	SeriesDescription string `json:"series_description"`
	UserId            uint   `json:"user_id"`
}

type FileAttachmentId struct {
	Id uint `json:"id"`
}

type PageQuery struct {
	Offset uint `json:"offset"`
	Limit  uint `json:"limit"`
}
