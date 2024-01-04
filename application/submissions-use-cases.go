package application

import (
	"github.com/upb-code-labs/tests-microservice/domain/definitions"
	"github.com/upb-code-labs/tests-microservice/domain/dtos"
	"github.com/upb-code-labs/tests-microservice/domain/entities"
	"github.com/upb-code-labs/tests-microservice/infrastructure/static_files"
)

type SubmissionsUseCases struct{}

func (submissionsUseCases *SubmissionsUseCases) RunTests(
	submissionWork *entities.SubmissionWork,
	testsRunner definitions.LanguageTestsRunner,
) error {
	// Get the archives
	staticFilesManager := static_files.StaticFilesManager{}

	languageTemplateArchive, err := staticFilesManager.GetLanguageTemplateBytes(submissionWork.LanguageUUID)
	if err != nil {
		return err
	}

	testsArchive, err := staticFilesManager.GetArchiveBytes(&dtos.GetFileFromMicroserviceDTO{
		FileUUID: submissionWork.TestsFileUUID,
		FileType: "test",
	})
	if err != nil {
		return err
	}

	solutionArchive, err := staticFilesManager.GetArchiveBytes(&dtos.GetFileFromMicroserviceDTO{
		FileUUID: submissionWork.SubmissionFileUUID,
		FileType: "submission",
	})
	if err != nil {
		return err
	}

	// Save the archives in the FS
	err = testsRunner.SaveArchivesInFS(&dtos.TestArchivesDTO{
		SubmissionUUID:          submissionWork.SubmissionUUID,
		LanguageTemplateArchive: &languageTemplateArchive,
		SubmissionArchive:       &solutionArchive,
		TestsArchive:            &testsArchive,
	})
	if err != nil {
		return err
	}

	// "Merge" the archives
	err = testsRunner.MergeArchives(submissionWork.SubmissionUUID)
	if err != nil {
		return err
	}

	// Run the tests
	err = testsRunner.RunTests(submissionWork.SubmissionUUID)
	if err != nil {
		return err
	}

	return nil
}
