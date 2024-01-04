package implementations

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type ArchivesManagerImplementation struct{}

func (archivesManagerImplementation *ArchivesManagerImplementation) SaveArchiveInFS(archiveBytes *[]byte, destinationPath string) error {
	// Create the file
	file, err := os.Create(destinationPath)
	if err != nil {
		return err
	}

	defer file.Close()

	// Write the file
	_, err = file.Write(*archiveBytes)
	if err != nil {
		return err
	}

	// Reset the file pointer
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}

	// Save the file
	err = file.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (archivesManagerImplementation *ArchivesManagerImplementation) ExtractArchive(archivePath string, destinationPath string) error {
	context, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	cmd := exec.CommandContext(
		context,
		"unzip",
		archivePath,
		"-d",
		destinationPath,
	)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (archivesManagerImplementation *ArchivesManagerImplementation) MoveFilesFromDirectoryToDirectory(sourceDirectoryPath string, destinationDirectoryPath string) error {
	context, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	moveCommand := fmt.Sprintf(
		"mv %s %s",
		sourceDirectoryPath,
		destinationDirectoryPath,
	)

	cmd := exec.CommandContext(
		context,
		"sh",
		"-c",
		moveCommand,
	)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (archivesManagerImplementation *ArchivesManagerImplementation) DeleteArchive(archivePath string) error {
	err := os.Remove(archivePath)
	if err != nil {
		return err
	}

	return nil
}

func (archivesManagerImplementation *ArchivesManagerImplementation) DeleteDirectory(directoryPath string) error {
	err := os.RemoveAll(directoryPath)
	if err != nil {
		return err
	}

	return nil
}
