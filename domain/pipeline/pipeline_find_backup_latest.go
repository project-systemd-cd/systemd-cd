package pipeline

import (
	errorss "errors"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/unix"
)

// findBackupLatest implements iPipeline
func (p *pipeline) findBackupLatest() (backupPath string, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Find latest backup")
	logger.Logger().Debugf("* pipeline.Name = %v", p.ManifestMerged.Name)
	logger.Logger().Tracef("* pipeline = %+v", *p)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		var ErrNotFound *errors.ErrNotFound
		if err == nil || errorss.As(err, &ErrNotFound) {
			logger.Logger().Debugf("> backupPath = %v", backupPath)
			logger.Logger().Debug("END   - Find latest backup")
		} else {
			logger.Logger().Error("FAILED - Find latest backup")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	backupBasePath := p.service.PathBackupDir + p.ManifestMerged.Name + "/"
	err = unix.MkdirIfNotExist(backupBasePath)
	if err != nil {
		return "", err
	}
	s, err := unix.Ls(
		unix.ExecuteOption{},
		unix.LsOption{SortByDescendingTime: true},
		backupBasePath,
	)
	if err != nil {
		return "", err
	}
	if len(s) == 0 {
		err = &errors.ErrNotFound{Object: "backup"}
		return "", err
	}

	backupPath = backupBasePath + s[0] + "/"

	return backupPath, nil
}
