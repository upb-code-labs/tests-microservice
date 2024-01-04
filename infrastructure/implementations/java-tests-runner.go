package implementations

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/upb-code-labs/tests-microservice/domain/dtos"
	"github.com/upb-code-labs/tests-microservice/infrastructure"
)

type JavaTestsRunner struct{}

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
		fmt.Sprintf("%s/template.zip", path),
	)
	if err != nil {
		return err
	}

	err = javaTestsRunner.saveArchiveInFS(
		dto.TestsArchive,
		fmt.Sprintf("%s/tests.zip", path),
	)
	if err != nil {
		return err
	}

	err = javaTestsRunner.saveArchiveInFS(
		dto.SubmissionArchive,
		fmt.Sprintf("%s/submission.zip", path),
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
		fmt.Sprintf("%s/template.zip", submissionPathPrefix),
		fmt.Sprintf("%s/template", submissionPathPrefix),
	)
	if err != nil {
		log.Println("Error while extracting template archive", err)
		return err
	}

	err = archivesManager.ExtractArchive(
		fmt.Sprintf("%s/tests.zip", submissionPathPrefix),
		fmt.Sprintf("%s/tests", submissionPathPrefix),
	)
	if err != nil {
		log.Println("Error while extracting tests archive", err)
		return err
	}

	err = archivesManager.ExtractArchive(
		fmt.Sprintf("%s/submission.zip", submissionPathPrefix),
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
		fmt.Sprintf("%s/template.zip", submissionPathPrefix),
	)
	if err != nil {
		log.Println("Error while deleting template archive", err)
		return err
	}

	err = archivesManager.DeleteArchive(
		fmt.Sprintf("%s/tests.zip", submissionPathPrefix),
	)
	if err != nil {
		log.Println("Error while deleting tests archive", err)
		return err
	}

	err = archivesManager.DeleteArchive(
		fmt.Sprintf("%s/submission.zip", submissionPathPrefix),
	)
	if err != nil {
		log.Println("Error while deleting submission archive", err)
		return err
	}

	return nil
}

func (javaTestsRunner *JavaTestsRunner) RunTests(submissionUUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the submission path
	submissionPath := fmt.Sprintf(
		"%s/%s/template/java_template",
		infrastructure.GetEnvironment().TestsExecutionDirectory,
		submissionUUID,
	)

	// Run the tests
	testAndBuildCommand := fmt.Sprintf(
		"cd %s && mvn clean test",
		submissionPath,
	)

	cmd := exec.CommandContext(
		ctx,
		"sh",
		"-c",
		testAndBuildCommand,
	)

	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Error while running tests", err)
		return err
	}

	fmt.Println("Tests output", string(out))

	// Delete the submission directory
	err = javaTestsRunner.deleteSubmissionDirectory(submissionUUID)
	if err != nil {
		return err
	}

	return nil
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
