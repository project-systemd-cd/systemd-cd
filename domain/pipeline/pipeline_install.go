package pipeline

import (
	"strings"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/systemd"
	"systemd-cd/domain/unix"
)

func (p pipeline) install() ([]systemd.UnitService, error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: p}}))

	var pathWorkingDir *string
	if p.ManifestMerged.Opt != nil {
		pathOptDir := p.service.PathOptDir + p.ManifestMerged.Name + "/"
		err := unix.MkdirIfNotExist(pathOptDir)
		if err != nil {
			return nil, err
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
				return nil, err
			}
		}
	}

	if p.ManifestMerged.Binaries != nil && len(*p.ManifestMerged.Binaries) != 0 {
		pathBinDir := p.service.PathBinDir + p.ManifestMerged.Name + "/"
		err := unix.MkdirIfNotExist(pathBinDir)
		if err != nil {
			return nil, err
		}
		for _, binary := range *p.ManifestMerged.Binaries {
			pathBinFile := pathBinDir + strings.TrimPrefix(binary, "./")

			// Copy binary files
			err := unix.Cp(
				unix.ExecuteOption{},
				unix.CpOption{Force: true},
				string(p.RepositoryLocal.Path)+"/"+binary,
				pathBinFile,
			)
			if err != nil {
				logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
				return nil, err
			}
		}
	}

	systemdServices := []systemd.UnitService{}
	if p.ManifestMerged.SystemdOptions != nil && len(p.ManifestMerged.SystemdOptions) != 0 {
		pathEtcDir := p.service.PathEtcDir + p.ManifestMerged.Name + "/"
		err := unix.MkdirIfNotExist(pathEtcDir)
		if err != nil {
			return nil, err
		}

		for _, service := range p.ManifestMerged.SystemdOptions {
			args := service.Args
			if strings.TrimSpace(args) != "" && !strings.HasPrefix(args, " ") {
				args = " " + args
			}
			if service.Etc != nil {
				// Copy or create etc files and add to cli options
				for _, etc := range service.Etc {
					etcFilePath := pathEtcDir + etc.Target
					if etc.Content != nil {
						// Create etc file
						err := unix.WriteFile(etcFilePath, []byte(*etc.Content))
						if err != nil {
							logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
							return nil, err
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
							return nil, err
						}
					}
					// Add to cli options
					args += " " + etc.Option + " " + pathEtcDir + etc.Target
				}
			}

			pathEnvFile := p.service.PathSystemdUnitEnvFileDir + service.Name
			env := map[string]string{}
			if service.EnvVars != nil {
				// Set environment variables
				for _, e := range service.EnvVars {
					env[e.Name] = e.Value
				}
			}

			execStart := strings.TrimPrefix(service.ExecuteCommand, "./")
			if p.ManifestMerged.Binaries != nil && len(*p.ManifestMerged.Binaries) != 0 {
				for _, binary := range *p.ManifestMerged.Binaries {
					pathBinDir := p.service.PathBinDir + p.ManifestMerged.Name + "/"
					pathBinFile := pathBinDir + strings.TrimPrefix(binary, "./")

					// If binary file name equals execute command, change to absolute path
					if strings.Split(strings.TrimPrefix(service.ExecuteCommand, "./"), " ")[0] ==
						strings.TrimPrefix(binary, "./") {
						// Cut out cli args
						args := strings.TrimPrefix(
							execStart,
							strings.TrimPrefix(binary, "./"),
						)
						// Create command for `ExecStart` in systemd unit
						execStart = pathBinFile + args
						break
					}
				}
			}

			unitType := systemd.UnitTypeSimple
			s, err := p.service.Systemd.NewService(
				service.Name,
				systemd.UnitFileService{
					Unit: systemd.UnitDirective{
						Description:   service.Description,
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
				return systemdServices, err
			}
			systemdServices = append(systemdServices, s)
		}
	}

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Value: systemdServices}}))
	return systemdServices, nil
}
