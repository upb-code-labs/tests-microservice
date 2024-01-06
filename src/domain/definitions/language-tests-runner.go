package definitions

import (
	"github.com/upb-code-labs/tests-microservice/src/domain/dtos"
)

type LanguageTestsRunner interface {
	RunTests(submissionUUID string) (*dtos.TestResultDTO, error)
	SaveArchivesInFS(dto *dtos.TestArchivesDTO) error
	MergeArchives(submissionUUID string) error
}
