package definitions

type ArchivesManager interface {
	SaveArchiveInFS(archiveBytes *[]byte, destinationPath string) error
	ExtractArchive(archivePath string, destinationPath string) error
	MoveFilesFromDirectoryToDirectory(sourceDirectoryPath string, destinationDirectoryPath string) error
	DeleteArchive(archivePath string) error
	DeleteDirectory(directoryPath string) error
}
