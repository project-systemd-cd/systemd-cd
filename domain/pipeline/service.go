package pipeline

import (
	"strings"
	"systemd-cd/domain/git"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/systemd"
	"systemd-cd/domain/unix"
)

type IPipelineService interface {
	// Resgister pipeline.
	// If pipeline with same name already exists, replace it.
	NewPipeline(ServiceManifestLocal) (IPipeline, error)

	FindPipelines() ([]PipelineMetadata, error)
	FindPipelineByName(name string) (PipelineMetadata, error)
	FindSystemdServiceByName(name string) (systemd.IService, error)
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

func NewService(repo IRepository, git git.IService, systemd systemd.IService, d Directories) (p IPipelineService, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Instantiate pipeline service")
	logger.Logger().Debugf("< dirSrc = %v", d.Src)
	logger.Logger().Debugf("< dirBinary = %v", d.Binary)
	logger.Logger().Debugf("< dirEtc = %v", d.Etc)
	logger.Logger().Debugf("< dirOpt = %v", d.Opt)
	logger.Logger().Debugf("< dirSystemdUnitFile = %v", d.SystemdUnitFile)
	logger.Logger().Debugf("< dirSystemdUnitEnvFile = %v", d.SystemdUnitEnvFile)
	logger.Logger().Debugf("< dirBackup = %v", d.Backup)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debug("END   - Instantiate pipeline service")
		} else {
			logger.Logger().Error("FAILED - Instantiate pipeline service")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	for _, d := range []*string{
		&d.Src, &d.Binary, &d.Etc, &d.Opt,
		&d.SystemdUnitFile, &d.SystemdUnitEnvFile, &d.Backup,
	} {
		if !strings.HasSuffix(*d, "/") {
			// add trailing slash
			*d += "/"
		}
		// Create directory
		err = unix.MkdirIfNotExist(*d)
		if err != nil {
			return pipelineService{}, err
		}
	}

	p = &pipelineService{
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
