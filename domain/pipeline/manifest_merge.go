package pipeline

import "systemd-cd/domain/logger"

func (remote ServiceManifestRemote) merge(remoteUrl string, local ServiceManifestLocal) (ServiceManifestMerged, error) {
	// Merge to local manifest
	var manifestRemoteSystemdOptions []SystemdOptionMerged = nil
	for _, s := range remote.SystemdOptions {
		description := remoteUrl
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
		Name:            remote.Name,
		GitTargetBranch: local.GitTargetBranch,
		GitTagRegex:     local.GitTagRegex,
		TestCommands:    remote.TestCommands,
		BuildCommands:   remote.BuildCommands,
		Opt:             remote.Opt,
		Binaries:        remote.Binaries,
		SystemdOptions:  manifestRemoteSystemdOptions,
	}
	manifestMerged.Name = local.Name
	if local.TestCommands != nil {
		manifestMerged.TestCommands = local.TestCommands
	}
	if local.BuildCommands != nil {
		manifestMerged.BuildCommands = local.BuildCommands
	}
	if local.Opt != nil {
		manifestMerged.Opt = *local.Opt
	}
	if local.Binaries != nil {
		manifestMerged.Binaries = local.Binaries
	}
	var systemdOptions []SystemdOptionMerged = nil
	for _, s := range local.SystemdOptions {
		description := remoteUrl
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
	err := manifestMerged.Validate()
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return ServiceManifestMerged{}, err
	}

	return manifestMerged, nil
}
