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

func (p pipeline) getRemoteManifest() (m ServiceManifestRemote, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Get manifest in git repository")
	logger.Logger().Debugf("* pipeline.Name = %v", p.ManifestMerged.Name)
	logger.Logger().Tracef("* pipeline = %+v", p)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debugf("> manifestRemote.Name = %v", m.Name)
			logger.Logger().Tracef("> manifestRemote = %+v", m)
			logger.Logger().Debug("END   - Get manifest in git repository")
		} else {
			logger.Logger().Error("FAILED - Get manifest in git repository")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	//* NOTE: No error if file not found

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
	b := &bytes.Buffer{}
	err = unix.ReadFile(manifestFilePath, b)
	if err != nil && !os.IsNotExist(err) {
		return ServiceManifestRemote{}, err
	}
	fileExists := !os.IsNotExist(err)
	if fileExists {
		err = toml.Decode(b, &m)
		if err != nil {
			return ServiceManifestRemote{}, err
		}
	}

	return m, nil
}
