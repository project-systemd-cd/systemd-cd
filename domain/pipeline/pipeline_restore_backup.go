package pipeline

import (
	"systemd-cd/domain/logger"
	"systemd-cd/domain/unix"
)

func (p pipeline) restoreBackup(o restoreBackupOptions) (err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Restore pipeline files from backup")
	logger.Logger().Debugf("* pipeline.Name = %v", p.ManifestMerged.Name)
	logger.Logger().Tracef("* pipeline = %+v", p)
	logger.Logger().Debugf("< option = %+v", o)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debug("END   - Restore pipeline files from backup")
		} else {
			logger.Logger().Error("FAILED - Restore pipeline files from backup")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

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

	if p.ManifestMerged.Binaries != nil {
		// Restore binary file
		err = unix.Cp(
			unix.ExecuteOption{},
			unix.CpOption{Force: true},
			backupPath+"bin/*",
			p.service.PathBinDir+p.ManifestMerged.Name+"/*",
		)
		if err != nil {
			return err
		}
	}

	if p.ManifestMerged.SystemdServiceOptions != nil && len(p.ManifestMerged.SystemdServiceOptions) != 0 {
		// Restore systemd unit file
		err = unix.Cp(
			unix.ExecuteOption{},
			unix.CpOption{Force: true},
			backupPath+"systemd/*",
			p.service.PathSystemdUnitFileDir,
		)
		if err != nil {
			return err
		}

		// Restore env file
		err = unix.Cp(
			unix.ExecuteOption{},
			unix.CpOption{Force: true},
			backupPath+"env/*",
			p.service.PathSystemdUnitEnvFileDir,
		)
		if err != nil {
			return err
		}

		for _, s := range p.ManifestMerged.SystemdServiceOptions {
			if len(s.Etc) != 0 {
				// Restore etc files
				err = unix.Cp(
					unix.ExecuteOption{},
					unix.CpOption{Recursive: true, Force: true},
					backupPath+"etc/*",
					p.service.PathEtcDir,
				)
				if err != nil {
					return err
				}
				break
			}
		}

		for _, s := range p.ManifestMerged.SystemdServiceOptions {
			if len(s.Opt) != 0 {
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
				break
			}
		}
	}

	return nil
}
