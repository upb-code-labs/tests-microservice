package utils

import (
	"errors"

	"github.com/upb-code-labs/tests-microservice/domain/definitions"
	"github.com/upb-code-labs/tests-microservice/infrastructure/implementations"
)

func GetTestRunnerByLanguageUUID(languageUUID string) (runner definitions.LanguageTestsRunner, err error) {
	JAVA_LANGUAGE_UUID := "487034c9-441c-4fb9-b0f3-8f4dd6176532"

	switch languageUUID {
	case JAVA_LANGUAGE_UUID:
		return &implementations.JavaTestsRunner{}, nil
	}

	return nil, errors.New("language not found")
}
