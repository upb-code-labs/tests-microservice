package dtos

import (
	"regexp"
)

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
	SubmissionUUID string `json:"submission_uuid"`
	TestsPassed    bool   `json:"tests_passed"`
	TestsOutput    string `json:"tests_output"`
}

func (dto *TestResultDTO) ToSubmissionStatusUpdateDTO(status string) *SubmissionStatusUpdateDTO {
	return &SubmissionStatusUpdateDTO{
		SubmissionUUID:   dto.SubmissionUUID,
		SubmissionStatus: status,
		TestsPassed:      dto.TestsPassed,
		TestsOutput:      dto.TestsOutput,
	}
}

type SubmissionStatusUpdateDTO struct {
	SubmissionUUID   string `json:"submission_uuid"`
	SubmissionStatus string `json:"submission_status"`
	TestsPassed      bool   `json:"tests_passed"`
	TestsOutput      string `json:"tests_output"`
}
