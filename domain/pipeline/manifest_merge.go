package pipeline

import "systemd-cd/domain/logger"

func (remote ServiceManifestRemote) merge(remoteUrl string, local ServiceManifestLocal) (m ServiceManifestMerged, err error) {
	logger.Logger().Debug("START - Merge local manifest to remote manifest")
	logger.Logger().Debugf("< localManifest.Name = %v", local.Name)
	logger.Logger().Debugf("< remoteUrl = %v", remoteUrl)
	defer func() {
		if err == nil {
			logger.Logger().Debug("END   - Merge local manifest to remote manifest")
		} else {
			logger.Logger().Error("FAILED - Merge local manifest to remote manifest")
			logger.Logger().Error(err)
		}
	}()

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
			Opt:            s.Opt,
			Port:           s.Port,
		})
	}
	m = ServiceManifestMerged{
		Name:            remote.Name,
		GitTargetBranch: local.GitTargetBranch,
		GitTagRegex:     local.GitTagRegex,
		TestCommands:    remote.TestCommands,
		BuildCommands:   remote.BuildCommands,
		Binaries:        remote.Binaries,
		SystemdOptions:  manifestRemoteSystemdOptions,
	}
	m.Name = local.Name
	if local.TestCommands != nil {
		m.TestCommands = local.TestCommands
	}
	if local.BuildCommands != nil {
		m.BuildCommands = local.BuildCommands
	}
	if local.Binaries != nil {
		m.Binaries = local.Binaries
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
			Opt:            s.Opt,
			Port:           s.Port,
		})
	}
	if systemdOptions != nil {
		m.SystemdOptions = systemdOptions
	}

	// Validate manifest
	err = m.Validate()
	if err != nil {
		return ServiceManifestMerged{}, err
	}

	return m, nil
}
