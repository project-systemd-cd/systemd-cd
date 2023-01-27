package pipeline

import (
	"strconv"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/unix"
	"time"
)

func (p pipeline) backupInstalled() error {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: p}}))

	// Create directory for backup
	// e.g.
	// /var/backups/systemd-cd/<name>/<unix-time>_<commit-id>/
	backupPath := p.service.PathBackupDir + p.ManifestMerged.Name + "/" + strconv.FormatInt(time.Now().Unix(), 10) + "_" + p.GetCommitRef() + "/"
	err := unix.MkdirIfNotExist(backupPath)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	// Backup systemd unit file
	// e.g.
	// `cp /usr/local/lib/systemd/system/<name>.service /var/backups/systemd-cd/<name>/<unix-time>_<commit-id>/unit.service`
	err = unix.Mv(
		unix.ExecuteOption{},
		unix.MvOption{},
		p.service.PathSystemdUnitFileDir+p.ManifestMerged.Name+".service",
		backupPath+"unit.service",
	)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))

	}

	// Backup env file
	// e.g.
	// `cp /usr/local/systemd-cd/etc/default/<name> /var/backups/systemd-cd/<name>/<unix-time>_<commit-id>/env`
	err = unix.Mv(
		unix.ExecuteOption{},
		unix.MvOption{},
		p.service.PathSystemdUnitEnvFileDir+p.ManifestMerged.Name,
		backupPath+"env",
	)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))

	}

	// TODO: check condition with old manifest.
	if p.ManifestMerged.Binaries != nil && len(*p.ManifestMerged.Binaries) != 0 {
		// Backup binary
		// e.g.
		// `cp /usr/local/systemd-cd/bin/<name>/* /var/backups/systemd-cd/<name>/<unix-time>_<commit-id>/bin/`
		err = unix.MkdirIfNotExist(backupPath + "bin/")
		if err != nil {
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return err
		}
		err = unix.Mv(
			unix.ExecuteOption{},
			unix.MvOption{},
			p.service.PathBinDir+p.ManifestMerged.Name+"/*",
			backupPath+"bin/",
		)
		if err != nil {
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))

		}
	}

	// TODO: check condition with old manifest.
	if len(p.ManifestMerged.Opt) != 0 {
		// Backup opt
		// e.g.
		// `cp /usr/local/systemd-cd/opt/<name>/* /var/backups/systemd-cd/<name>/<unix-time>_<commit-id>/opt/`
		err = unix.MkdirIfNotExist(backupPath + "opt/")
		if err != nil {
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return err
		}
		err = unix.Mv(
			unix.ExecuteOption{},
			unix.MvOption{},
			p.service.PathOptDir+p.ManifestMerged.Name+"/*",
			backupPath+"opt/",
		)
		if err != nil {
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))

		}
	}

	logger.Logger().Trace("Finished")
	return nil
}
