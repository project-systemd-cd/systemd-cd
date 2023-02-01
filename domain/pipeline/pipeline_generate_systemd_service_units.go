package pipeline

import (
	"strings"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/systemd"
)

func (p *pipeline) generateSystemdServiceUnits() (units []systemdUnit) {
	logger.Logger().Debug("START - Generate systemd service units")
	logger.Logger().Debugf("* pipeline.Name = %v", p.ManifestMerged.Name)
	logger.Logger().Tracef("* pipeline = %+v", *p)
	defer func() {
		for i, su := range units {
			logger.Logger().Debugf("> units[%d].Name = %v", i, su.Name)
			logger.Logger().Tracef("> units[%d] = %+v", i, su)
		}
		logger.Logger().Debug("END   - Generate systemd service units")
	}()

	unitType := systemd.UnitTypeSimple

	if p.ManifestMerged.SystemdOptions != nil && len(p.ManifestMerged.SystemdOptions) != 0 {
		for _, service := range p.ManifestMerged.SystemdOptions {
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

			args := service.Args
			if strings.TrimSpace(args) != "" && !strings.HasPrefix(args, " ") {
				args = " " + args
			}

			pathEnvFile := p.service.PathSystemdUnitEnvFileDir + service.Name
			env := map[string]string{}
			if service.EnvVars != nil {
				// Set environment variables
				for _, e := range service.EnvVars {
					env[e.Name] = e.Value
				}
			}

			argsEtc := ""
			if service.Etc != nil {
				pathEtcDir := p.service.PathEtcDir + service.Name + "/"

				// Copy or create etc files and add to cli options
				for _, etc := range service.Etc {
					// Add to cli options
					argsEtc += " " + etc.Option + " " + pathEtcDir + etc.Target
				}
			}

			var pathWorkingDir *string
			if service.Opt != nil {
				pathOptDir := p.service.PathOptDir + service.Name + "/"
				pathWorkingDir = &pathOptDir
			}

			units = append(units, systemdUnit{
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
						ExecStart:        execStart + args + argsEtc,
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
			})
		}
	}

	return units
}
