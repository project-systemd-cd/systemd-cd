package pipeline

import (
	"systemd-cd/domain/errors"
	"systemd-cd/domain/unix"
)

// findBackupLatest implements iPipeline
func (p *pipeline) findBackupLatest() (backupPath string, err error) {
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
