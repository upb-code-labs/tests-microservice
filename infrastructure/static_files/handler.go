package static_files

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/upb-code-labs/tests-microservice/domain/dtos"
	"github.com/upb-code-labs/tests-microservice/infrastructure"
)

type StaticFilesManager struct{}

func (staticFilesManager *StaticFilesManager) GetArchiveBytes(dto *dtos.GetFileFromMicroserviceDTO) ([]byte, error) {
	endpoint := fmt.Sprintf(
		"%s/archives/download",
		infrastructure.GetEnvironment().StaticFilesMicroserviceAddress,
	)

	// Create request payload from the dto
	body, err := json.Marshal(dto)
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}

	// Create request
	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return []byte{}, err
	}

	request.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return []byte{}, err
	}

	defer response.Body.Close()

	// Handle error
	if response.StatusCode != http.StatusOK {
		return []byte{}, errors.New(
			"there was an error while trying to get the archive from the static files microservice",
		)
	}

	// Read response
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	// Return response
	return responseBody, nil
}

func (staticFilesManager *StaticFilesManager) GetLanguageTemplateBytes(languageUUID string) ([]byte, error) {
	endpoint := fmt.Sprintf(
		"%s/templates/%s",
		infrastructure.GetEnvironment().StaticFilesMicroserviceAddress,
		languageUUID,
	)

	// Create request
	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return []byte{}, err
	}

	// Send the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return []byte{}, err
	}

	defer response.Body.Close()

	// Handle error
	if response.StatusCode != http.StatusOK {
		return []byte{}, errors.New(
			"there was an error while trying to get the template from the static files microservice",
		)
	}

	// Read response
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	// Return response
	return responseBody, nil
}
