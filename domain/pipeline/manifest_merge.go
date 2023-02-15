package pipeline

import "systemd-cd/domain/logger"

func (remote ServiceManifestRemote) merge(remoteUrl string, local ServiceManifestLocal) (m ServiceManifestMerged, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Merge local manifest to remote manifest")
	logger.Logger().Tracef("* remoteManifest = %+v", remote)
	logger.Logger().Debugf("< localManifest.Name = %v", local.Name)
	logger.Logger().Debugf("< remoteUrl = %v", remoteUrl)
	logger.Logger().Tracef("< localManifest = %+v", local)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Tracef("> manifestMerged = %+v", m)
			logger.Logger().Debug("END   - Merge local manifest to remote manifest")
		} else {
			logger.Logger().Error("FAILED - Merge local manifest to remote manifest")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	// Merge to local manifest
	var manifestRemoteSystemdOptions []SystemdServiceOptionMerged = nil
	for _, s := range remote.SystemdOptions {
		description := remoteUrl
		if s.Description != nil && *s.Description != "" {
			description = *s.Description
		}
		manifestRemoteSystemdOptions = append(manifestRemoteSystemdOptions, SystemdServiceOptionMerged{
			Name:         s.Name,
			Description:  description,
			ExecStartPre: s.ExecStartPre,
			ExecStart:    s.ExecStart,
			Args:         s.Args,
			EnvVars:      s.EnvVars,
			Etc:          s.Etc,
			Opt:          s.Opt,
			Port:         s.Port,
		})
	}
	m = ServiceManifestMerged{
		Name:                  remote.Name,
		GitTargetBranch:       local.GitTargetBranch,
		GitTagRegex:           local.GitTagRegex,
		TestCommands:          remote.TestCommands,
		BuildCommands:         remote.BuildCommands,
		Binaries:              remote.Binaries,
		SystemdServiceOptions: manifestRemoteSystemdOptions,
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
	var systemdOptions []SystemdServiceOptionMerged = nil
	for _, s := range local.SystemdOptions {
		description := remoteUrl
		if s.Description != nil && *s.Description != "" {
			description = *s.Description
		}
		systemdOptions = append(systemdOptions, SystemdServiceOptionMerged{
			Name:         s.Name,
			Description:  description,
			ExecStartPre: s.ExecStartPre,
			ExecStart:    s.ExecStart,
			Args:         s.Args,
			EnvVars:      s.EnvVars,
			Etc:          s.Etc,
			Opt:          s.Opt,
			Port:         s.Port,
		})
	}
	if systemdOptions != nil {
		m.SystemdServiceOptions = systemdOptions
	}

	// Validate manifest
	err = m.Validate()
	if err != nil {
		return ServiceManifestMerged{}, err
	}

	return m, nil
}
