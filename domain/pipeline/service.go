package pipeline

import (
	"strings"
	"systemd-cd/domain/git"
	"systemd-cd/domain/systemd"
	"systemd-cd/domain/unix"
)

type IPipelineService interface {
	NewPipeline(ServiceManifestLocal) (IPipeline, error)

	FindPipelines() ([]PipelineMetadata, error)
	FindPipeline(name string) (PipelineMetadata, error)
	FindSystemdService(name string) (systemd.IService, error)
}

type Directories struct {
	Src                string
	Binary             string
	Etc                string
	Opt                string
	SystemdUnitFile    string
	SystemdUnitEnvFile string
	Backup             string
}

func NewService(repo IRepository, git git.IService, systemd systemd.IService, d Directories) (IPipelineService, error) {
	for _, d := range []*string{
		&d.Src, &d.Binary, &d.Etc, &d.Opt,
		&d.SystemdUnitFile, &d.SystemdUnitEnvFile, &d.Backup,
	} {
		if !strings.HasSuffix(*d, "/") {
			// add trailing slash
			*d += "/"
		}
		// Create directory
		err := unix.MkdirIfNotExist(*d)
		if err != nil {
			return pipelineService{}, err
		}
	}

	p := &pipelineService{
		repo, git, systemd,
		d.Src, d.Binary, d.Etc, d.Opt,
		d.SystemdUnitFile, d.SystemdUnitEnvFile, d.Backup,
	}

	return p, nil
}

type pipelineService struct {
	repo    IRepository
	Git     git.IService
	Systemd systemd.IService

	PathSrcDir                string
	PathBinDir                string
	PathEtcDir                string
	PathOptDir                string
	PathSystemdUnitFileDir    string
	PathSystemdUnitEnvFileDir string
	PathBackupDir             string
}

// FindPipeline implements IPipelineService
func (s pipelineService) FindPipeline(name string) (PipelineMetadata, error) {
	m, err := s.repo.FindPipeline(name)
	if err != nil {
		return PipelineMetadata{}, err
	}

	return m, nil
}

// FindPipelines implements IPipelineService
func (s pipelineService) FindPipelines() ([]PipelineMetadata, error) {
	metadatas, err := s.repo.FindPipelines()
	if err != nil {
		return nil, err
	}

	return metadatas, nil
}

// FindSystemdService implements IPipelineService
func (s pipelineService) FindSystemdService(name string) (systemd.IService, error) {
	panic("unimplemented")
}
