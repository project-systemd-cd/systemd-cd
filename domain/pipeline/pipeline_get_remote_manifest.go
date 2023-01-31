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

func (p pipeline) getRemoteManifest() (ServiceManifestRemote, error) {
	//* NOTE: No error if file not found
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
		return ServiceManifestRemote{}, err
	}
	fileExists := !os.IsNotExist(err)
	if fileExists {
		err = toml.Decode(b, manifestRemote)
		if err != nil {
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return ServiceManifestRemote{}, err
		}
	}

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Value: *manifestRemote}}))
	return *manifestRemote, nil
}
