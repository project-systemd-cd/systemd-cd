package pipeline

import (
	"systemd-cd/domain/unix"
)

func (p pipeline) restoreBackup(o restoreBackupOptions) (err error) {
	// Find backup
	backupPath := ""
	if o.CommidId != nil {
		backupPath, err = p.findBackupByCommitId(*o.CommidId)
	} else {
		backupPath, err = p.findBackupLatest()
	}
	if err != nil {
		return err
	}

	// Restore systemd unit file
	err = unix.Cp(
		unix.ExecuteOption{},
		unix.CpOption{Force: true},
		backupPath+"unit.service",
		p.service.PathSystemdUnitFileDir+p.ManifestMerged.Name+".service",
	)
	if err != nil {
		return err
	}

	// Restore env file
	err = unix.Cp(
		unix.ExecuteOption{},
		unix.CpOption{Force: true},
		backupPath+"env",
		p.service.PathSystemdUnitEnvFileDir+p.ManifestMerged.Name,
	)
	if err != nil {
		return err
	}

	if p.ManifestMerged.Binary != nil {
		// Restore binary file
		err = unix.Cp(
			unix.ExecuteOption{},
			unix.CpOption{Force: true},
			backupPath+"binary",
			p.service.PathBinDir+p.ManifestMerged.Name+"/"+*p.ManifestMerged.Binary,
		)
		if err != nil {
			return err
		}
	}

	if p.ManifestMerged.Opt != nil {
		// Restore opt files
		err = unix.Cp(
			unix.ExecuteOption{},
			unix.CpOption{Recursive: true, Force: true},
			backupPath+"opt/*",
			p.service.PathOptDir,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
