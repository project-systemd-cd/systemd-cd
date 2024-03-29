package pipeline

import (
	"strconv"
	"strings"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/unix"
	"time"
)

func (p pipeline) backupInstalled() (err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Backup pipeline files")
	logger.Logger().Debugf("* pipeline.Name = %v", p.ManifestMerged.Name)
	logger.Logger().Tracef("* pipeline = %+v", p)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debug("END   - Backup pipeline files")
		} else {
			logger.Logger().Error("FAILED - Backup pipeline files")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	// Create directory for backup
	// e.g.
	// /var/backups/systemd-cd/<name>/<unix-time>_<commit-id>/
	backupPath := p.service.PathBackupDir + p.ManifestMerged.Name + "/" + strconv.FormatInt(time.Now().Unix(), 10) + "_" + p.GetCommitRef() + "/"
	err = unix.MkdirIfNotExist(backupPath)
	if err != nil {
		return err
	}

	if p.ManifestMerged.SystemdServiceOptions != nil && len(p.ManifestMerged.SystemdServiceOptions) != 0 {
		err = unix.MkdirIfNotExist(backupPath+"systemd/", backupPath+"env/")
		if err != nil {
			return err
		}
		for _, s := range p.ManifestMerged.SystemdServiceOptions {
			// Backup systemd unit file
			// e.g.
			// `cp /usr/local/lib/systemd/system/<unit_name>.service /var/backups/systemd-cd/<name>/<unix-time>_<commit-id>/systemd/<unit_name>.service`
			err = unix.Mv(
				unix.ExecuteOption{},
				unix.MvOption{},
				p.service.PathSystemdUnitFileDir+s.Name+".service",
				backupPath+"systemd/",
			)
			if err != nil {
				return err
			}

			// Backup env file
			// e.g.
			// `cp /usr/local/systemd-cd/etc/default/<unit_name> /var/backups/systemd-cd/<name>/<unix-time>_<commit-id>/env/<unit_name>`
			err = unix.Mv(
				unix.ExecuteOption{},
				unix.MvOption{},
				p.service.PathSystemdUnitEnvFileDir+s.Name,
				backupPath+"env/",
			)
			if err != nil {
				return err
			}

			if len(s.Etc) != 0 {
				// Backup etc file
				// e.g.
				// `cp /usr/local/systemd-cd/etc/<unit_name>/* /var/backups/systemd-cd/<name>/<unix-time>_<commit-id>/etc/<unit_name>/`
				err = unix.MkdirIfNotExist(backupPath + "etc/" + s.Name)
				if err != nil {
					return err
				}
				for _, etc := range s.Etc {
					if !strings.HasPrefix(strings.TrimPrefix(etc.Target, "./"), ".") {
						err = unix.Mv(
							unix.ExecuteOption{},
							unix.MvOption{},
							p.service.PathEtcDir+s.Name+"/*",
							backupPath+"etc/"+s.Name+"/",
						)
						if err != nil {
							return err
						}
						break
					}
				}
				for _, etc := range s.Etc {
					// Backup hidden files
					filename := strings.TrimPrefix(etc.Target, "./")
					if strings.HasPrefix(etc.Target, ".") {
						err = unix.Mv(
							unix.ExecuteOption{},
							unix.MvOption{},
							p.service.PathEtcDir+s.Name+"/"+filename,
							backupPath+"etc/"+s.Name+"/",
						)
						if err != nil {
							return err
						}
					}
				}
			}

			if len(s.Opt) != 0 {
				// Backup opt files
				// e.g.
				// `cp /usr/local/systemd-cd/opt/<unit_name>/* /var/backups/systemd-cd/<name>/<unix-time>_<commit-id>/opt/<unit_name>/`
				err = unix.MkdirIfNotExist(backupPath + "opt/" + s.Name)
				if err != nil {
					return err
				}
				for _, filename := range s.Opt {
					if !strings.HasPrefix(strings.TrimPrefix(filename, "./"), ".") {
						err = unix.Mv(
							unix.ExecuteOption{},
							unix.MvOption{},
							p.service.PathOptDir+s.Name+"/*",
							backupPath+"opt/"+s.Name+"/",
						)
						if err != nil {
							return err
						}
						break
					}
				}
				for _, filename := range s.Opt {
					// Backup hidden files
					filename = strings.TrimPrefix(filename, "./")
					if strings.HasPrefix(filename, ".") {
						err = unix.Mv(
							unix.ExecuteOption{},
							unix.MvOption{},
							p.service.PathOptDir+s.Name+"/"+filename,
							backupPath+"opt/"+s.Name+"/",
						)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	if p.ManifestMerged.Binaries != nil && len(*p.ManifestMerged.Binaries) != 0 {
		// Backup binary files
		// e.g.
		// `cp /usr/local/systemd-cd/bin/<name>/* /var/backups/systemd-cd/<name>/<unix-time>_<commit-id>/bin/`
		err = unix.MkdirIfNotExist(backupPath + "bin/")
		if err != nil {
			return err
		}
		err = unix.Mv(
			unix.ExecuteOption{},
			unix.MvOption{},
			p.service.PathBinDir+p.ManifestMerged.Name+"/*",
			backupPath+"bin/",
		)
		if err != nil {
			return err
		}
	}

	return nil
}
