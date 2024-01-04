package dtos

import "regexp"

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

type ReplaceRegexDTO struct {
	Regexp      regexp.Regexp
	Replacement string
}

type TestResultDTO struct {
	SubmissionUUID string
	TestsPassed    bool
	TestsOutput    string
}
