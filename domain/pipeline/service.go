package pipeline

import (
	"strings"
	"systemd-cd/domain/git"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/systemd"
	"systemd-cd/domain/unix"
)

type iPipelineService interface {
	NewPipeline(ServiceManifestLocal) (iPipeline, error)
}

type Directories struct {
	Var                string
	Src                string
	Binary             string
	Etc                string
	Opt                string
	SystemdUnitFile    string
	SystemdUnitEnvFile string
	Backup             string
}

func NewService(git git.IService, systemd systemd.IService, d Directories) (iPipelineService, error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "directories", Value: d}}))

	for _, d := range []*string{
		&d.Var, &d.Src, &d.Binary, &d.Etc, &d.Opt,
		&d.SystemdUnitFile, &d.SystemdUnitEnvFile, &d.Backup,
	} {
		if !strings.HasSuffix(*d, "/") {
			// add trailing slash
			*d += "/"
		}
		// Create directory
		err := unix.MkdirIfNotExist(*d)
		if err != nil {
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return pipelineService{}, err
		}
	}

	p := &pipelineService{
		git, systemd,
		d.Var, d.Src, d.Binary, d.Etc, d.Opt,
		d.SystemdUnitFile, d.SystemdUnitEnvFile, d.Backup,
	}

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Value: *p}}))
	return p, nil
}

type pipelineService struct {
	Git     git.IService
	Systemd systemd.IService

	PathVarDir                string
	PathSrcDir                string
	PathBinDir                string
	PathEtcDir                string
	PathOptDir                string
	PathSystemdUnitFileDir    string
	PathSystemdUnitEnvFileDir string
	PathBackupDir             string
}