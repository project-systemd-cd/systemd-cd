package pipeline

import (
	"errors"
	"systemd-cd/domain/unix"
)

// findBackupLatest implements iPipeline
func (p *pipeline) findBackupLatest() (string, error) {
	backupBasePath := p.service.PathBackupDir + p.ManifestMerged.Name + "/"
	s, err := unix.Ls(
		unix.ExecuteOption{},
		unix.LsOption{SortByDescendingTime: true},
		backupBasePath,
	)
	if err != nil {
		return "", err
	}
	if len(s) == 0 {
		return "", errors.New("no backups")
	}

	backupPath := backupBasePath + s[0] + "/"

	return backupPath, nil
}
