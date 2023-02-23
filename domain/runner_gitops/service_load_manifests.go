package runner_gitops

import (
	"bytes"
	"os"
	"strings"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/toml"
	"systemd-cd/domain/unix"
)

func (s *service) loadManifests() (manifests []pipeline.ServiceManifestLocal, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Instantiate gitops runner service")
	logger.Logger().Debugf("> service.repository = %+v", *s.repository)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debugf("> manifests = %+v", manifests)
			logger.Logger().Debug("END   - Instantiate gitops runner service")
		} else {
			logger.Logger().Error("FAILED - Instantiate gitops runner service")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	_, b, _, err := unix.Execute(
		unix.ExecuteOption{WorkingDirectory: (*string)(&s.repository.Path)},
		"/usr/bin/find", ".", "-type", "f", "-name", "'*.toml'",
	)
	if err != nil {
		return nil, err
	}
	files := strings.Split(b.String(), "\n")
	for _, filename := range files {
		if filename == "" {
			continue
		}
		// Read file
		manifestLocal := new(pipeline.ServiceManifestLocal)
		b := &bytes.Buffer{}
		path := (string)(s.repository.Path) + "/" + strings.TrimPrefix(filename, "./")
		logger.Logger().Debugf("Load manifest %s", path)
		err = unix.ReadFile(path, b)
		if err != nil {
			return nil, err
		}
		fileExists := !os.IsNotExist(err)
		if fileExists {
			err = toml.Decode(b, manifestLocal)
			if err != nil {
				return nil, err
			}
			manifests = append(manifests, *manifestLocal)
		}
	}
	return manifests, nil
}
