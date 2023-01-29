package pipeline

import (
	"bytes"
	"os"
	"strings"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/toml"
	"systemd-cd/domain/unix"
)

const defaultManifestFileName = ".systemd-cd.yaml"

func (p pipeline) loadManifest() (ServiceManifestMerged, error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: p}}))

	// Get paths
	repositoryPath := string(p.RepositoryLocal.Path)
	if !strings.HasSuffix(repositoryPath, "/") {
		repositoryPath += "/"
	}
	manifestFilePath := repositoryPath + defaultManifestFileName
	if p.ManifestLocal.GitManifestFile != nil {
		manifestFilePath = repositoryPath + *p.ManifestLocal.GitManifestFile
	}

	// Read file
	manifestRemote := new(ServiceManifestRemote)
	b := &bytes.Buffer{}
	err := unix.ReadFile(manifestFilePath, b)
	if err != nil && !os.IsNotExist(err) {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return ServiceManifestMerged{}, err
	}
	if !os.IsNotExist(err) {
		// If manifest file exists, unmarshal to struct
		err = toml.Decode(b, manifestRemote)
		if err != nil {
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return ServiceManifestMerged{}, err
		}
	}

	// Merge to local manifest
	var manifestRemoteSystemdOptions []SystemdOptionMerged = nil
	for _, s := range manifestRemote.SystemdOptions {
		description := p.RepositoryLocal.RemoteUrl
		if s.Description != nil && *s.Description != "" {
			description = *s.Description
		}
		manifestRemoteSystemdOptions = append(manifestRemoteSystemdOptions, SystemdOptionMerged{
			Name:           s.Name,
			Description:    description,
			ExecuteCommand: s.ExecuteCommand,
			Args:           s.Args,
			EnvVars:        s.EnvVars,
			Etc:            s.Etc,
			Port:           s.Port,
		})
	}
	manifestMerged := ServiceManifestMerged{
		Name:           manifestRemote.Name,
		TestCommands:   manifestRemote.TestCommands,
		BuildCommands:  manifestRemote.BuildCommands,
		Opt:            manifestRemote.Opt,
		Binaries:       manifestRemote.Binaries,
		SystemdOptions: manifestRemoteSystemdOptions,
	}
	manifestMerged.Name = p.ManifestLocal.Name
	if p.ManifestLocal.TestCommands != nil {
		manifestMerged.TestCommands = p.ManifestLocal.TestCommands
	}
	if p.ManifestLocal.BuildCommands != nil {
		manifestMerged.BuildCommands = p.ManifestLocal.BuildCommands
	}
	if p.ManifestLocal.Opt != nil {
		manifestMerged.Opt = *p.ManifestLocal.Opt
	}
	if p.ManifestLocal.Binaries != nil {
		manifestMerged.Binaries = p.ManifestLocal.Binaries
	}
	var systemdOptions []SystemdOptionMerged = nil
	for _, s := range p.ManifestLocal.SystemdOptions {
		description := p.RepositoryLocal.RemoteUrl
		if s.Description != nil && *s.Description != "" {
			description = *s.Description
		}
		systemdOptions = append(systemdOptions, SystemdOptionMerged{
			Name:           s.Name,
			Description:    description,
			ExecuteCommand: s.ExecuteCommand,
			Args:           s.Args,
			EnvVars:        s.EnvVars,
			Etc:            s.Etc,
			Port:           s.Port,
		})
	}
	if systemdOptions != nil {
		manifestMerged.SystemdOptions = systemdOptions
	}

	// Validate manifest
	err = manifestMerged.Validate()
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return ServiceManifestMerged{}, err
	}

	return manifestMerged, nil
}
