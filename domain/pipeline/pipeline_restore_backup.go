package pipeline

import (
	"systemd-cd/domain/logger"
	"systemd-cd/domain/unix"
)

func (p pipeline) restoreBackup(o restoreBackupOptions) (err error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: p}, {Value: o}}))

	// Find backup
	backupPath := ""
	if o.CommidId != nil {
		backupPath, err = p.findBackupByCommitId(*o.CommidId)
	} else {
		backupPath, err = p.findBackupLatest()
	}
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
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
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return err
		}
	}

	if p.ManifestMerged.SystemdOptions != nil && len(p.ManifestMerged.SystemdOptions) != 0 {
		// Restore systemd unit file
		err = unix.Cp(
			unix.ExecuteOption{},
			unix.CpOption{Force: true},
			backupPath+"systemd/*",
			p.service.PathSystemdUnitFileDir,
		)
		if err != nil {
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
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
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return err
		}
	}

	for _, s := range p.ManifestMerged.SystemdOptions {
		if s.Opt != nil && len(s.Opt) != 0 {
			// Restore opt files
			err = unix.Cp(
				unix.ExecuteOption{},
				unix.CpOption{Recursive: true, Force: true},
				backupPath+"opt/*",
				p.service.PathOptDir,
			)
			if err != nil {
				logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
				return err
			}
			break
		}
	}

	logger.Logger().Trace("Finished")
	return nil
}
