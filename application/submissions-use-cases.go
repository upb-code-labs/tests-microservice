package application

import (
	"github.com/upb-code-labs/tests-microservice/domain/definitions"
	"github.com/upb-code-labs/tests-microservice/domain/dtos"
	"github.com/upb-code-labs/tests-microservice/domain/entities"
	"github.com/upb-code-labs/tests-microservice/infrastructure/static_files"
	"github.com/upb-code-labs/tests-microservice/utils"
)

type SubmissionsUseCases struct{}

func (submissionsUseCases *SubmissionsUseCases) RunTests(
	submissionWork *entities.SubmissionWork,
	testsRunner definitions.LanguageTestsRunner,
) dtos.TestResultDTO {
	// Get the archives
	staticFilesManager := static_files.StaticFilesManager{}

	languageTemplateArchive, err := staticFilesManager.GetLanguageTemplateBytes(submissionWork.LanguageUUID)
	if err != nil {
		return *utils.GetTestResultDTOFromErrorMessage(
			submissionWork.SubmissionUUID,
			"[ERROR] We couldn't get the programming language archive to run the tests",
		)
	}

	testsArchive, err := staticFilesManager.GetArchiveBytes(&dtos.GetFileFromMicroserviceDTO{
		FileUUID: submissionWork.TestsFileUUID,
		FileType: "test",
	})
	if err != nil {
		return *utils.GetTestResultDTOFromErrorMessage(
			submissionWork.SubmissionUUID,
			"[ERROR] We couldn't get the tests archive to run the tests",
		)
	}

	solutionArchive, err := staticFilesManager.GetArchiveBytes(&dtos.GetFileFromMicroserviceDTO{
		FileUUID: submissionWork.SubmissionFileUUID,
		FileType: "submission",
	})
	if err != nil {
		return *utils.GetTestResultDTOFromErrorMessage(
			submissionWork.SubmissionUUID,
			"[ERROR] We couldn't get your submission archive to run the tests",
		)
	}

	// Save the archives in the FS
	err = testsRunner.SaveArchivesInFS(&dtos.TestArchivesDTO{
		SubmissionUUID:          submissionWork.SubmissionUUID,
		LanguageTemplateArchive: &languageTemplateArchive,
		SubmissionArchive:       &solutionArchive,
		TestsArchive:            &testsArchive,
	})
	if err != nil {
		return *utils.GetTestResultDTOFromErrorMessage(
			submissionWork.SubmissionUUID,
			"[ERROR] We couldn't save the archives in the file system to run the tests",
		)
	}

	// "Merge" the archives
	err = testsRunner.MergeArchives(submissionWork.SubmissionUUID)
	if err != nil {
		return *utils.GetTestResultDTOFromErrorMessage(
			submissionWork.SubmissionUUID,
			"[ERROR] We couldn't prepare the archives to run the tests",
		)
	}

	// Run the tests
	result, _ := testsRunner.RunTests(submissionWork.SubmissionUUID)
	return *result
}
