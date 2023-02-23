package main

import (
	"bytes"
	"os"
	"strings"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/toml"
	"systemd-cd/domain/unix"
)

func loadManifests(paths []string, recursive bool) (manifests []pipeline.ServiceManifestLocal, err error) {
	logger.Logger().Info("-----------------------------------------------------------")
	logger.Logger().Info("START - Load manifests")
	logger.Logger().Infof("< paths = %+v", paths)
	logger.Logger().Info("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Info("-----------------------------------------------------------")
		if err == nil {
			for i, sml := range manifests {
				logger.Logger().Infof("> manifests[%d].Name = %v", i, sml.Name)
				logger.Logger().Tracef("> manifests[%d] = %+v", i, sml)
			}
			logger.Logger().Info("END   - Load manifests")
		} else {
			logger.Logger().Error("FAILED - Load manifests")
			logger.Logger().Error(err)
		}
		logger.Logger().Info("-----------------------------------------------------------")
	}()

	for _, path := range paths {
		if recursive {
			_, b, _, err := unix.Execute(unix.ExecuteOption{}, "/usr/bin/find", path, "-type", "f", "-name", "'*.toml'")
			if err != nil {
				return nil, err
			}
			files := strings.Split(b.String(), "\n")
			for _, filename := range files {
				if filename == "" {
					continue
				}
				logger.Logger().Infof("Load manifest file \"%s\"", filename)
				sml, err := loadManifest(filename)
				if err != nil {
					return nil, err
				}
				manifests = append(manifests, sml)
			}
		} else {
			if !strings.HasSuffix(path, ".toml") {
				if !strings.HasSuffix(path, "/") {
					path += "/"
				}
				path += "*.toml"
			}
			filenames, err := unix.Ls(unix.ExecuteOption{}, unix.LsOption{DirTrailiingSlash: true}, path)
			if err != nil {
				return nil, err
			}
			for _, filename := range filenames {
				if filename == "" {
					continue
				}
				sml, err := loadManifest(filename)
				if err != nil {
					return nil, err
				}
				manifests = append(manifests, sml)
			}
		}
	}

	return manifests, nil
}

func loadManifest(path string) (pipeline.ServiceManifestLocal, error) {
	// Read file
	manifestLocal := new(pipeline.ServiceManifestLocal)
	b := &bytes.Buffer{}
	err := unix.ReadFile(path, b)
	if err != nil && !os.IsNotExist(err) {
		return pipeline.ServiceManifestLocal{}, err
	}
	fileExists := !os.IsNotExist(err)
	if fileExists {
		err = toml.Decode(b, manifestLocal)
		if err != nil {
			return pipeline.ServiceManifestLocal{}, err
		}
	}

	return *manifestLocal, nil
}
