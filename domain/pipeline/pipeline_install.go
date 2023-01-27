package pipeline

import (
	"strings"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/systemd"
	"systemd-cd/domain/unix"
)

func (p pipeline) install() (systemd.UnitService, error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: p}}))

	pathEnvFile := p.service.PathSystemdUnitEnvFileDir + p.ManifestMerged.Name

	var pathWorkingDir *string
	if p.ManifestMerged.Opt != nil {
		pathOptDir := p.service.PathOptDir + p.ManifestMerged.Name + "/"
		err := unix.MkdirIfNotExist(pathOptDir)
		if err != nil {
			return systemd.UnitService{}, err
		}
		pathWorkingDir = &pathOptDir

		// Copy opt files
		for _, src := range p.ManifestMerged.Opt {
			// Copy opt file
			err := unix.Cp(
				unix.ExecuteOption{WorkingDirectory: (*string)(&p.RepositoryLocal.Path)},
				unix.CpOption{Recursive: true, Parents: true, Force: true},
				src,
				pathOptDir+src,
			)
			if err != nil {
				logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
				return systemd.UnitService{}, err
			}
		}
	}

	args := p.ManifestMerged.Args
	if strings.TrimSpace(args) != "" && !strings.HasPrefix(args, " ") {
		args = " " + args
	}
	if p.ManifestLocal.Etc != nil {
		pathEtcDir := p.service.PathEtcDir + p.ManifestMerged.Name + "/"
		err := unix.MkdirIfNotExist(pathEtcDir)
		if err != nil {
			return systemd.UnitService{}, err
		}

		// Copy or create etc files and add to cli options
		for _, etc := range *p.ManifestLocal.Etc {
			etcFilePath := pathEtcDir + etc.Target
			if etc.Content != nil {
				// Create etc file
				err := unix.WriteFile(etcFilePath, []byte(*etc.Content))
				if err != nil {
					logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
					return systemd.UnitService{}, err
				}
			} else {
				// Copy etc file
				err := unix.Cp(
					unix.ExecuteOption{WorkingDirectory: (*string)(&p.RepositoryLocal.Path)},
					unix.CpOption{Recursive: true, Force: true},
					etc.Target,
					etcFilePath,
				)
				if err != nil {
					logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
					return systemd.UnitService{}, err
				}
			}
			// Add to cli options
			args += " " + etc.Option + " " + pathEtcDir + etc.Target
		}
	}

	env := map[string]string{}
	if p.ManifestMerged.EnvVars != nil {
		// Set environment variables
		for _, e := range p.ManifestMerged.EnvVars {
			env[e.Name] = e.Value
		}
	}

	execStart := strings.TrimPrefix(p.ManifestMerged.ExecuteCommand, "./")
	if p.ManifestMerged.Binaries != nil && len(*p.ManifestMerged.Binaries) != 0 {
		pathBinDir := p.service.PathBinDir + p.ManifestMerged.Name + "/"
		err := unix.MkdirIfNotExist(pathBinDir)
		if err != nil {
			return systemd.UnitService{}, err
		}
		for _, binary := range *p.ManifestMerged.Binaries {
			pathBinFile := pathBinDir + strings.TrimPrefix(binary, "./")

			// Copy binary files
			err = unix.Cp(
				unix.ExecuteOption{},
				unix.CpOption{Force: true},
				string(p.RepositoryLocal.Path)+"/"+binary,
				pathBinFile,
			)
			if err != nil {
				logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
				return systemd.UnitService{}, err
			}

			// If binary file name equals execute command, change to absolute path
			if strings.Split(strings.TrimPrefix(p.ManifestMerged.ExecuteCommand, "./"), " ")[0] ==
				strings.TrimPrefix(binary, "./") {
				// Cut out cli args
				args := strings.TrimPrefix(
					execStart,
					strings.TrimPrefix(binary, "./"),
				)
				// Create command for `ExecStart` in systemd unit
				logger.Logger().Debugf("Debug:\n\tpathBinFile: %v\n\targs: %v", pathBinFile, args)
				execStart = pathBinFile + args
			}
		}
	}

	unitType := systemd.UnitTypeSimple
	service, err := p.service.Systemd.NewService(
		p.ManifestMerged.Name,
		systemd.UnitFileService{
			Unit: systemd.UnitDirective{
				Description:   p.ManifestMerged.Description,
				Documentation: p.RepositoryLocal.RemoteUrl,
				After:         nil,
				Requires:      nil,
				Wants:         nil,
				Conflicts:     nil,
			},
			Service: systemd.ServiceDirective{
				Type:             &unitType,
				WorkingDirectory: pathWorkingDir,
				EnvironmentFile:  &pathEnvFile,
				ExecStart:        execStart + args,
				ExecStop:         nil,
				ExecReload:       nil,
				Restart:          nil,
				RemainAfterExit:  nil,
			},
			Install: systemd.InstallDirective{
				Alias:           nil,
				RequiredBy:      nil,
				WantedBy:        []string{"multi-user.target"},
				Also:            nil,
				DefaultInstance: nil,
			},
		},
		env,
	)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return service, err
	}

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Value: service}}))
	return service, err
}
