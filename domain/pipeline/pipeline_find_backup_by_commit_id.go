package pipeline

import (
	"strings"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/unix"
)

// findBackupByCommitId implements iPipeline
func (p *pipeline) findBackupByCommitId(commitId string) (backupPath string, err error) {
	logger.Logger().Debug("START - Find backup by commit id")
	logger.Logger().Debugf("< pipeline.Name = %v", p.ManifestMerged.Name)
	defer func() {
		if err == nil {
			logger.Logger().Debugf("> backupPath = %v", backupPath)
			logger.Logger().Debug("END   - Find backup by commit id")
		} else {
			logger.Logger().Error("FAILED - Find backup by commit id")
			logger.Logger().Error(err)
		}
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
		err = &errors.ErrNotFound{Object: "backup", IdName: "version", Id: commitId}
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
		return "", err
	}

	return backupPath, nil
}
