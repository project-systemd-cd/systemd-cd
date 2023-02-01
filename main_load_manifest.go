package main

import (
	"bytes"
	"os"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/toml"
	"systemd-cd/domain/unix"
)

func loadManifest(path string) (pipeline.ServiceManifestLocal, error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "path", Value: path}}))

	// Read file
	manifestLocal := new(pipeline.ServiceManifestLocal)
	b := &bytes.Buffer{}
	err := unix.ReadFile(path, b)
	if err != nil && !os.IsNotExist(err) {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return pipeline.ServiceManifestLocal{}, err
	}
	fileExists := !os.IsNotExist(err)
	if fileExists {
		err = toml.Decode(b, manifestLocal)
		if err != nil {
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return pipeline.ServiceManifestLocal{}, err
		}
	}

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Value: *manifestLocal}}))
	return *manifestLocal, nil
}
