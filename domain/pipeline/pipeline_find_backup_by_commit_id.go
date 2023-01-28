package pipeline

import (
	"errors"
	"fmt"
	"strings"
	"systemd-cd/domain/unix"
)

// findBackupByCommitId implements iPipeline
func (p *pipeline) findBackupByCommitId(commitId string) (string, error) {
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
		return "", fmt.Errorf("backup of version '%s' not found", commitId)
	}

	return backupPath, nil
}
