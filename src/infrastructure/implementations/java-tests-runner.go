package implementations

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/upb-code-labs/tests-microservice/src/domain/dtos"
	"github.com/upb-code-labs/tests-microservice/src/infrastructure"
)

type JavaTestsRunner struct{}

var templateArchivePathTemplate = "%s/template.zip"
var testsArchivePathTemplate = "%s/tests.zip"
var submissionArchivePathTemplate = "%s/submission.zip"

// SaveArchivesInFS saves the archives needed to run the tests in the file system
func (javaTestsRunner *JavaTestsRunner) SaveArchivesInFS(dto *dtos.TestArchivesDTO) error {
	// Ensure the directory doesn't exist
	path := fmt.Sprintf(
		"%s/%s",
		infrastructure.GetEnvironment().TestsExecutionDirectory,
		dto.SubmissionUUID,
	)

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return errors.New("tests execution directory already exists, skipping")
	}

	// Create the directory
	err := os.Mkdir(path, 0700)
	if err != nil {
		return err
	}

	// Save the archives
	err = javaTestsRunner.saveArchiveInFS(
		dto.LanguageTemplateArchive,
		fmt.Sprintf(templateArchivePathTemplate, path),
	)
	if err != nil {
		return err
	}

	err = javaTestsRunner.saveArchiveInFS(
		dto.TestsArchive,
		fmt.Sprintf(testsArchivePathTemplate, path),
	)
	if err != nil {
		return err
	}

	err = javaTestsRunner.saveArchiveInFS(
		dto.SubmissionArchive,
		fmt.Sprintf(submissionArchivePathTemplate, path),
	)
	if err != nil {
		return err
	}

	return nil
}

func (javaTestsRunner *JavaTestsRunner) saveArchiveInFS(fileBytes *[]byte, path string) error {
	archivesManager := ArchivesManagerImplementation{}
	err := archivesManager.SaveArchiveInFS(fileBytes, path)
	return err
}

// MergeArchives merges the content of the template archive, teacher's tests archive and student's
// submission archive into a single directory that will be used to run the tests
func (javaTestsRunner *JavaTestsRunner) MergeArchives(submissionUUID string) error {
	// Unzip the archives
	err := javaTestsRunner.unzipArchives(submissionUUID)
	if err != nil {
		return err
	}

	// Delete the archives
	err = javaTestsRunner.deleteArchives(submissionUUID)
	if err != nil {
		return err
	}

	// Merge the archives
	archivesManager := ArchivesManagerImplementation{}

	err = archivesManager.MoveFilesFromDirectoryToDirectory(
		fmt.Sprintf("%s/%s/tests/*/src/test/java", infrastructure.GetEnvironment().TestsExecutionDirectory, submissionUUID),
		fmt.Sprintf("%s/%s/template/java_template/src/test/", infrastructure.GetEnvironment().TestsExecutionDirectory, submissionUUID),
	)
	if err != nil {
		return err
	}

	err = archivesManager.MoveFilesFromDirectoryToDirectory(
		fmt.Sprintf("%s/%s/submission/*/src/main/java", infrastructure.GetEnvironment().TestsExecutionDirectory, submissionUUID),
		fmt.Sprintf("%s/%s/template/java_template/src/main/", infrastructure.GetEnvironment().TestsExecutionDirectory, submissionUUID),
	)
	if err != nil {
		return err
	}

	return nil
}

func (javaTestsRunner *JavaTestsRunner) unzipArchives(submissionUUID string) error {
	submissionPathPrefix := fmt.Sprintf(
		"%s/%s",
		infrastructure.GetEnvironment().TestsExecutionDirectory,
		submissionUUID,
	)

	archivesManager := ArchivesManagerImplementation{}

	err := archivesManager.ExtractArchive(
		fmt.Sprintf(templateArchivePathTemplate, submissionPathPrefix),
		fmt.Sprintf("%s/template", submissionPathPrefix),
	)
	if err != nil {
		log.Println("Error while extracting template archive", err)
		return err
	}

	err = archivesManager.ExtractArchive(
		fmt.Sprintf(testsArchivePathTemplate, submissionPathPrefix),
		fmt.Sprintf("%s/tests", submissionPathPrefix),
	)
	if err != nil {
		log.Println("Error while extracting tests archive", err)
		return err
	}

	err = archivesManager.ExtractArchive(
		fmt.Sprintf(submissionArchivePathTemplate, submissionPathPrefix),
		fmt.Sprintf("%s/submission", submissionPathPrefix),
	)
	if err != nil {
		log.Println("Error while extracting submission archive", err)
		return err
	}

	return nil
}

func (javaTestsRunner *JavaTestsRunner) deleteArchives(submissionUUID string) error {
	submissionPathPrefix := fmt.Sprintf(
		"%s/%s",
		infrastructure.GetEnvironment().TestsExecutionDirectory,
		submissionUUID,
	)

	archivesManager := ArchivesManagerImplementation{}

	err := archivesManager.DeleteArchive(
		fmt.Sprintf(templateArchivePathTemplate, submissionPathPrefix),
	)
	if err != nil {
		log.Println("Error while deleting template archive", err)
		return err
	}

	err = archivesManager.DeleteArchive(
		fmt.Sprintf(testsArchivePathTemplate, submissionPathPrefix),
	)
	if err != nil {
		log.Println("Error while deleting tests archive", err)
		return err
	}

	err = archivesManager.DeleteArchive(
		fmt.Sprintf(submissionArchivePathTemplate, submissionPathPrefix),
	)
	if err != nil {
		log.Println("Error while deleting submission archive", err)
		return err
	}

	return nil
}

// RunTests runs the tests and returns the result
func (javaTestsRunner *JavaTestsRunner) RunTests(submissionUUID string) (dto *dtos.TestResultDTO, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Delete the submission directory at the end
	defer javaTestsRunner.deleteSubmissionDirectory(submissionUUID)

	// Get the submission path
	submissionPath := fmt.Sprintf(
		"%s/%s/template/java_template",
		infrastructure.GetEnvironment().TestsExecutionDirectory,
		submissionUUID,
	)

	// Prepare the command
	testAndBuildCommand := fmt.Sprintf(
		"cd %s && timeout 1m mvn clean test",
		submissionPath,
	)

	cmd := exec.CommandContext(
		ctx,
		"sh",
		"-c",
		testAndBuildCommand,
	)

	// Run the command
	out, err := cmd.CombinedOutput()

	// Parse errors (if any)
	if err != nil {
		errorLines := javaTestsRunner.getErrorLinesFromOutput(string(out))
		errorLines = javaTestsRunner.sanitizeConsoleTextLines(errorLines)

		return &dtos.TestResultDTO{
			SubmissionUUID: submissionUUID,
			TestsPassed:    false,
			TestsOutput:    strings.Join(errorLines, "\n"),
		}, nil
	}

	// Parse success lines
	successLines := javaTestsRunner.getResultLinesFromSuccessOutput(string(out))
	successLines = javaTestsRunner.sanitizeConsoleTextLines(successLines)

	return &dtos.TestResultDTO{
		SubmissionUUID: submissionUUID,
		TestsPassed:    true,
		TestsOutput:    strings.Join(successLines, "\n"),
	}, nil
}

// getErrorLinesFromOutput returns the lines starting with the `[ERROR]` prefix from the output
func (javaTestsRunner *JavaTestsRunner) getErrorLinesFromOutput(output string) []string {
	errorRegex := regexp.MustCompile(`(?m)^\[ERROR\].*$`)
	errorLines := errorRegex.FindAllString(output, -1)
	return errorLines
}

// sanitizeConsoleTextLines removes the lines that are not relevant for the user or can contain
// sensitive information
func (javaTestsRunner *JavaTestsRunner) sanitizeConsoleTextLines(textLines []string) []string {
	sanitizedTextLines := []string{}

	regExpToReplace := []dtos.ReplaceRegexDTO{
		// Remove possible path to the tests execution directory
		{
			Regexp:      *regexp.MustCompile(`\/[a-zA-Z0-9_\-]+(?:\/[a-zA-Z0-9_\-]+)*\/template`),
			Replacement: "****/****",
		},
		// Remove lines starting with [WARNING]
		{
			Regexp:      *regexp.MustCompile(`(?m)^\[WARNING\].*$`),
			Replacement: "",
		},
	}

	for _, errorLine := range textLines {
		for _, regExp := range regExpToReplace {
			errorLine = regExp.Regexp.ReplaceAllString(errorLine, regExp.Replacement)
		}

		if len(errorLine) > 0 {
			sanitizedTextLines = append(sanitizedTextLines, errorLine)
		}
	}

	return sanitizedTextLines
}

// getResultLinesFromSuccessOutput returns the lines of the output starting from the tests results
// summary / header line until the end
func (javaTestsRunner *JavaTestsRunner) getResultLinesFromSuccessOutput(output string) []string {
	// Header line
	successRegex := regexp.MustCompile(`\[INFO\] Tests run: \d+, Failures: \d+, Errors: \d+, Skipped: \d+`)

	// Find the header line and trim the text from it (including the header line) to the end
	headerLineIndex := successRegex.FindStringIndex(output)
	if headerLineIndex == nil {
		return []string{}
	}

	output = output[headerLineIndex[0]:]
	return strings.Split(output, "\n")
}

func (javaTestsRunner *JavaTestsRunner) deleteSubmissionDirectory(submissionUUID string) error {
	submissionPath := fmt.Sprintf(
		"%s/%s",
		infrastructure.GetEnvironment().TestsExecutionDirectory,
		submissionUUID,
	)

	err := os.RemoveAll(submissionPath)
	if err != nil {
		return err
	}

	return nil
}
