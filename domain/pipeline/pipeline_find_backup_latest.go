package pipeline

import (
	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/unix"
)

// findBackupLatest implements iPipeline
func (p *pipeline) findBackupLatest() (backupPath string, err error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: p}}))

	backupBasePath := p.service.PathBackupDir + p.ManifestMerged.Name + "/"
	err = unix.MkdirIfNotExist(backupBasePath)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return "", err
	}
	s, err := unix.Ls(
		unix.ExecuteOption{},
		unix.LsOption{SortByDescendingTime: true},
		backupBasePath,
	)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return "", err
	}
	if len(s) == 0 {
		err = &errors.ErrNotFound{Object: "backup"}
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return "", err
	}

	backupPath = backupBasePath + s[0] + "/"

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "backupBasePath", Value: backupBasePath}}))
	return backupPath, nil
}
