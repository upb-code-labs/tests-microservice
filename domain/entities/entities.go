package entities

type SubmissionWork struct {
	SubmissionUUID     string `json:"submission_uuid"`
	LanguageUUID       string `json:"language_uuid"`
	SubmissionFileUUID string `json:"submission_archive_uuid"`
	TestsFileUUID      string `json:"test_archive_uuid"`
}
