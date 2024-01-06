package utils

import (
	"errors"

	"github.com/upb-code-labs/tests-microservice/src/domain/definitions"
	"github.com/upb-code-labs/tests-microservice/src/domain/dtos"
	"github.com/upb-code-labs/tests-microservice/src/infrastructure/implementations"
)

func GetTestRunnerByLanguageUUID(languageUUID string) (runner definitions.LanguageTestsRunner, err error) {
	JAVA_LANGUAGE_UUID := "487034c9-441c-4fb9-b0f3-8f4dd6176532"

	switch languageUUID {
	case JAVA_LANGUAGE_UUID:
		return &implementations.JavaTestsRunner{}, nil
	}

	return nil, errors.New("language not found")
}

func GetTestResultDTOFromErrorMessage(
	submissionUUID string,
	errorMessage string,
) *dtos.TestResultDTO {
	return &dtos.TestResultDTO{
		SubmissionUUID: submissionUUID,
		TestsPassed:    false,
		TestsOutput:    errorMessage,
	}
}

func GetSubmissionStatusUpdateDTOFromErrorMessage(
	submissionUUID string,
	errorMessage string,
) *dtos.SubmissionStatusUpdateDTO {
	return &dtos.SubmissionStatusUpdateDTO{
		SubmissionUUID:   submissionUUID,
		TestsPassed:      false,
		TestsOutput:      errorMessage,
		SubmissionStatus: "ready",
	}
}
