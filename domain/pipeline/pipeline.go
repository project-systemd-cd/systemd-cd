package pipeline

import (
	"systemd-cd/domain/git"
	"systemd-cd/domain/systemd"
)

type IPipeline interface {
	GetName() string
	GetManifestLocal() ServiceManifestLocal

	// Execute test, execute build, and install.
	Init() error
	// If update exists, pull src and manifest, execute test, execute build, backup old files and install new version.
	// If fail to execute over systemd, restore from backup.
	Sync() error

	Uninstall() error

	GetStatus() status
	GetGitRemoteUrl() string
	GetGitTargetBranch() string
	GetGitTargetTagRegex() *string
	GetCommitRef() string
	GetCommitMessage() (string, error)
	GetStatusSystemdServices() ([]SystemdServiceWithStatus, error)

	GetJob(groupId string) ([]Job, error)
	GetJobs(QueryParamJob) ([][]Job, error)
}

type Path = string

type systemdUnit struct {
	Name     string
	UnitFile systemd.UnitFileService
	Env      map[string]string
}

type restoreBackupOptions struct {
	CommidId *string
}

type pipeline struct {
	ManifestLocal   ServiceManifestLocal
	ManifestMerged  ServiceManifestMerged
	RepositoryLocal *git.RepositoryLocal
	Status          status

	service *pipelineService
}

type SystemdServiceWithStatus struct {
	systemd.IUnitService
	Status systemd.Status
}
