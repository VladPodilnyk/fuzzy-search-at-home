package models

type DicomFileMetadata struct {
	PatientName string `json:"patient_name"`
	BirthDate   string `json:"birth_date"`
}

type PageQuery struct {
	Offset uint `json:"offset"`
	Limit  uint `json:"limit"`
}
