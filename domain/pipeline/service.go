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
	GetPipeline(name string) (IPipeline, error)

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
			return nil, err
		}
	}

	p = &pipelineService{
		repo, git, systemd,
		d.Src, d.Binary, d.Etc, d.Opt,
		d.SystemdUnitFile, d.SystemdUnitEnvFile, d.Backup,
	}

	go func() (err error) {
		defer func() {
			if err != nil {
				logger.Logger().Error("FAILED - Cancel pipeline job left behind")
				logger.Logger().Error(err)
			}
		}()
		pipelines, err := p.FindPipelines()
		if err != nil {
			return err
		}
		for _, pm := range pipelines {
			var p1 IPipeline
			p1, err = p.GetPipeline(pm.Name)
			if err != nil {
				return err
			}
			var jobGroups [][]Job
			jobGroups, err = p1.GetJobs(QueryParamJob{})
			if err != nil {
				return err
			}
			for _, jg := range jobGroups {
				for _, j := range jg {
					if j.Status == StatusJobInProgress || j.Status == StatusJobPending {
						// Cancel interrupted jobs
						err = jobInstance{Job: j}.Cancel(repo)
						if err != nil {
							return err
						}
					}
				}
			}
		}
		return nil
	}()

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
