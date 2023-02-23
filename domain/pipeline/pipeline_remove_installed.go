package pipeline

import (
	"systemd-cd/domain/logger"
	"systemd-cd/domain/unix"
)

func (p pipeline) removeInstalled() (err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Remove pipeline files")
	logger.Logger().Debugf("* pipeline.Name = %v", p.ManifestMerged.Name)
	logger.Logger().Tracef("* pipeline = %+v", p)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debug("END   - Remove pipeline files")
		} else {
			logger.Logger().Error("FAILED - Remove pipeline files")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	if p.ManifestMerged.SystemdServiceOptions != nil && len(p.ManifestMerged.SystemdServiceOptions) != 0 {
		for _, s := range p.ManifestMerged.SystemdServiceOptions {
			// Remove systemd unit file and env file
			// e.g.
			// `rm /usr/local/lib/systemd/system/<unit_name>.service /usr/local/systemd-cd/etc/default/<unit_name>`
			_, _, _, err = unix.Execute(
				unix.ExecuteOption{WantExitCodes: []int{1}},
				"rm",
				p.service.PathSystemdUnitFileDir+s.Name+".service",
				p.service.PathSystemdUnitEnvFileDir+s.Name,
			)
			if err != nil {
				return err
			}

			if len(s.Etc) != 0 {
				// Remove etc files
				// e.g.
				// `rm -r /usr/local/systemd-cd/etc/<unit_name>/`
				_, _, _, err = unix.Execute(
					unix.ExecuteOption{WantExitCodes: []int{1}},
					"rm", "-r",
					p.service.PathEtcDir+s.Name,
				)
				if err != nil {
					return err
				}
			}

			if len(s.Opt) != 0 {
				// Remove opt files
				// e.g.
				// `rm -r /usr/local/systemd-cd/opt/<unit_name>/`
				_, _, _, err = unix.Execute(
					unix.ExecuteOption{WantExitCodes: []int{1}},
					"rm", "-r",
					p.service.PathOptDir+s.Name,
				)
				if err != nil {
					return err
				}
			}
		}
	}

	if p.ManifestMerged.Binaries != nil && len(*p.ManifestMerged.Binaries) != 0 {
		// Remove binary files
		// e.g.
		// `rm /usr/local/systemd-cd/bin/<name>/`
		_, _, _, err = unix.Execute(
			unix.ExecuteOption{WantExitCodes: []int{1}},
			"rm", "-r",
			p.service.PathBinDir+p.ManifestMerged.Name,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
