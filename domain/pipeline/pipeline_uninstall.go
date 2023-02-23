package pipeline

import (
	errorss "errors"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"
)

func (p *pipeline) Uninstall() (err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Uninstall pipeline")
	logger.Logger().Debugf("* pipeline.Name = %v", p.ManifestLocal.Name)
	logger.Logger().Tracef("* pipeline = %+v", *p)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debug("END   - Uninstall pipeline")
		} else {
			logger.Logger().Error("FAILED - Uninstall pipeline")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	logger.Logger().Infof("Uninstall pipeline \"%s\"", p.ManifestMerged.Name)

	// Stop systemd service before backup
	systemdServices, err := p.getSystemdServices()
	for _, s := range systemdServices {
		logger.Logger().Debug("Stop systemd unit service \"%v\"", s.GetName())
		err = s.Disable(true)
		if err != nil {
			return err
		}
	}

	// Backup and remove installed files
	_, err = p.findBackupByCommitId(p.RepositoryLocal.RefCommitId)
	var ErrNotFound *errors.ErrNotFound
	notFound := errorss.As(err, &ErrNotFound)
	if notFound {
		err = p.backupInstalled()
	} else {
		err = p.removeInstalled()
	}

	return err
}
