package pipeline

import (
	"systemd-cd/domain/git"
	"systemd-cd/domain/systemd"
)

type iPipeline interface {
	// Execute test, execute build, and install.
	Init() error
	// If update exists, pull src and manifest, execute test, execute build, backup old files and install new version.
	// If fail to execute over systemd, restore from backup.
	Sync() error

	CheckUpdate() (updateExists bool, err error)
	GetStatus() Status
	GetCommitRef() string

	loadManifest() (ServiceManifestMerged, error)
	test() error
	build() error
	backupInstalled() error
	install() ([]systemd.UnitService, error)

	// Restore latest backup.
	// If `CommitId` specified, restore backup of specified version.
	restoreBackup(restoreBackupOptions) error

	findBackupByCommitId(commitId string) (Path, error)
	findBackupLatest() (Path, error)
}

type Path = string

type restoreBackupOptions struct {
	CommidId *string
}

type pipeline struct {
	ManifestLocal   ServiceManifestLocal
	ManifestMerged  ServiceManifestMerged
	RepositoryLocal *git.RepositoryLocal
	Status          Status

	service *pipelineService
}
