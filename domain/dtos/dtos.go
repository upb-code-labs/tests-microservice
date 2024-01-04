package dtos

type GetFileFromMicroserviceDTO struct {
	FileUUID string `json:"archive_uuid"`
	FileType string `json:"archive_type"`
}

type TestArchivesDTO struct {
	SubmissionUUID          string
	LanguageTemplateArchive *[]byte
	SubmissionArchive       *[]byte
	TestsArchive            *[]byte
}
