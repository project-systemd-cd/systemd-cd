package pipeline

import (
	"strings"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/unix"
)

// findBackupByCommitId implements iPipeline
func (p *pipeline) findBackupByCommitId(commitId string) (backupPath string, err error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: p}, {Name: "commitId", Value: commitId}}))

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
		err = &errors.ErrNotFound{Object: "backup", IdName: "version", Id: commitId}
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return "", err
	}

	backupPath = backupBasePath + s[0] + "/"
	// Search backups
	found := false
	for _, dir := range s {
		if strings.Contains(dir, commitId) {
			// found
			found = true
			backupPath = backupBasePath + dir + "/"
			break
		}
	}
	if !found {
		err = &errors.ErrNotFound{Object: "backup", IdName: "version", Id: commitId}
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return "", err
	}

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "backupBasePath", Value: backupBasePath}}))
	return backupPath, nil
}
