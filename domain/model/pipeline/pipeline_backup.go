package pipeline

import (
	"strconv"
	"systemd-cd/domain/model/logger"
	"systemd-cd/domain/model/unix"
	"time"
)

func (p pipeline) backupInstalled() error {
	logger.Logger().Tracef("Called:\n\tpipeline: %+v", p)

	// Create directory for backup
	// e.g.
	// /var/backups/systemd-cd/<name>/<unix-time>_<commit-id>/
	backupPath := p.service.PathBackupDir + p.ManifestMerged.Name + "/" + strconv.FormatInt(time.Now().Unix(), 10) + "_" + p.GetCommitRef() + "/"
	err := unix.MkdirIfNotExist(backupPath)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	// Backup systemd unit file
	// e.g.
	// `cp /usr/local/lib/systemd/system/<name>.service /var/backups/systemd-cd/<name>/<unix-time>_<commit-id>/unit.service`
	err = unix.Cp(
		unix.ExecuteOption{},
		unix.CpOption{},
		p.service.PathSystemdUnitFileDir+p.ManifestMerged.Name+".service",
		backupPath+"unit.service",
	)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	// Backup env file
	// e.g.
	// `cp /usr/local/systemd-cd/etc/default/<name> /var/backups/systemd-cd/<name>/<unix-time>_<commit-id>/env`
	err = unix.Cp(
		unix.ExecuteOption{},
		unix.CpOption{},
		p.service.PathSystemdUnitEnvFileDir+p.ManifestMerged.Name,
		backupPath+"env",
	)
	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v", err)
		return err
	}

	if p.ManifestMerged.Binary != nil {
		// Backup binary
		// e.g.
		// `cp /usr/local/systemd-cd/bin/<name>/<binary> /var/backups/systemd-cd/<name>/<unix-time>_<commit-id>/binary`
		err = unix.Cp(
			unix.ExecuteOption{},
			unix.CpOption{},
			p.service.PathBinDir+p.ManifestMerged.Name+"/"+*p.ManifestMerged.Binary,
			backupPath+"binary",
		)
		if err != nil {
			logger.Logger().Errorf("Error:\n\terr: %v", err)
			return err
		}
	}

	if p.ManifestMerged.Opt != nil {
		// Backup opt
		// e.g.
		// `cp /usr/local/systemd-cd/opt/<name>/* /var/backups/systemd-cd/<name>/<unix-time>_<commit-id>/opt/*`
		err = unix.Cp(
			unix.ExecuteOption{},
			unix.CpOption{
				Recursive: true,
			},
			p.service.PathOptDir+"*",
			backupPath+"opt/",
		)
		if err != nil {
			logger.Logger().Errorf("Error:\n\terr: %v", err)
			return err
		}
	}

	logger.Logger().Trace("Finished")
	return nil
}
