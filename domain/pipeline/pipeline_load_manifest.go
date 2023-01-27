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
	manifestInRepository := new(ServiceManifestRemote)
	b := &bytes.Buffer{}
	err := unix.ReadFile(manifestFilePath, b)
	if err != nil && !os.IsNotExist(err) {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return ServiceManifestMerged{}, err
	}
	if !os.IsNotExist(err) {
		// If manifest file exists, unmarshal to struct
		err = toml.Decode(b, manifestInRepository)
		if err != nil {
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return ServiceManifestMerged{}, err
		}
	}

	// Merge to local manifest
	manifestMerged := ServiceManifestMerged{
		Name:           manifestInRepository.Name,
		Description:    manifestInRepository.Description,
		Port:           manifestInRepository.Port,
		TestCommands:   manifestInRepository.TestCommands,
		BuildCommands:  manifestInRepository.BuildCommands,
		Opt:            manifestInRepository.Opt,
		Etc:            manifestInRepository.Etc,
		EnvVars:        p.ManifestLocal.EnvVars,
		Binaries:       manifestInRepository.Binaries,
		ExecuteCommand: manifestInRepository.ExecuteCommand,
		Args:           manifestInRepository.Args,
	}
	manifestMerged.Name = p.ManifestLocal.Name
	if p.ManifestLocal.Description != nil {
		manifestMerged.Description = *p.ManifestLocal.Description
	}
	if p.ManifestLocal.Port != nil {
		manifestMerged.Port = p.ManifestLocal.Port
	}
	if p.ManifestLocal.TestCommands != nil {
		manifestMerged.TestCommands = p.ManifestLocal.TestCommands
	}
	if p.ManifestLocal.BuildCommands != nil {
		manifestMerged.BuildCommands = p.ManifestLocal.BuildCommands
	}
	if p.ManifestLocal.Opt != nil {
		manifestMerged.Opt = *p.ManifestLocal.Opt
	}
	if p.ManifestLocal.Etc != nil {
		manifestMerged.Etc = *p.ManifestLocal.Etc
	}
	if p.ManifestLocal.Binaries != nil {
		manifestMerged.Binaries = p.ManifestLocal.Binaries
	}
	if p.ManifestLocal.ExecuteCommand != nil {
		manifestMerged.ExecuteCommand = *p.ManifestLocal.ExecuteCommand
	}
	if p.ManifestLocal.Args != nil {
		manifestMerged.Args = *p.ManifestLocal.Args
	}

	// Validate manifest
	err = manifestMerged.Validate()
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return ServiceManifestMerged{}, err
	}

	return manifestMerged, nil
}
