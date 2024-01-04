package definitions

import (
	"github.com/upb-code-labs/tests-microservice/domain/dtos"
)

type LanguageTestsRunner interface {
	RunTests(submissionUUID string) error
	SaveArchivesInFS(dto *dtos.TestArchivesDTO) error
	MergeArchives(submissionUUID string) error
}
