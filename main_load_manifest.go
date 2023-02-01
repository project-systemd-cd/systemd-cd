package main

import (
	"bytes"
	"os"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/toml"
	"systemd-cd/domain/unix"
)

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
