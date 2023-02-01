package pipeline

import (
	"strings"
	"systemd-cd/domain/systemd"
	"systemd-cd/domain/unix"
)

func (p pipeline) install() ([]systemd.UnitService, error) {
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
				return nil, err
			}
		}
	}

	systemdServices := []systemd.UnitService{}
	if p.ManifestMerged.SystemdOptions != nil && len(p.ManifestMerged.SystemdOptions) != 0 {
		for _, service := range p.ManifestMerged.SystemdOptions {
			if service.Etc != nil {
				pathEtcDir := p.service.PathEtcDir + service.Name + "/"
				err := unix.MkdirIfNotExist(pathEtcDir)
				if err != nil {
					return nil, err
				}

				// Copy or create etc files and add to cli options
				for _, etc := range service.Etc {
					etcFilePath := pathEtcDir + etc.Target
					if etc.Content != nil {
						// Create etc file
						err := unix.WriteFile(etcFilePath, []byte(*etc.Content))
						if err != nil {
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
							return nil, err
						}
					}
				}
			}

			if service.Opt != nil {
				pathOptDir := p.service.PathOptDir + service.Name + "/"
				err := unix.MkdirIfNotExist(pathOptDir)
				if err != nil {
					return nil, err
				}

				// Copy opt files
				for _, src := range service.Opt {
					// Copy opt file
					err := unix.Cp(
						unix.ExecuteOption{WorkingDirectory: (*string)(&p.RepositoryLocal.Path)},
						unix.CpOption{Recursive: true, Parents: true, Force: true},
						src,
						pathOptDir+src,
					)
					if err != nil {
						return nil, err
					}
				}
			}
		}

		for _, unit := range p.generateSystemdServiceUnits() {
			s, err := p.service.Systemd.NewService(unit.Name, unit.UnitFile, unit.Env)
			if err != nil {
				return nil, err
			}
			systemdServices = append(systemdServices, s)
		}
	}

	return systemdServices, nil
}
