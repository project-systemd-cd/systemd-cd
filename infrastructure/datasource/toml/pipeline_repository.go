package toml

import (
	"strings"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/unix"
)

func NewRepositoryPipeline(path string) (pipeline.IRepository, error) {
	err := unix.MkdirIfNotExist(path)
	if err != nil {
		return &rPipeline{}, err
	}
	if !strings.HasSuffix(path, "/") {
		// Add trailing slash
		path += "/"
	}

	err = unix.MkdirIfNotExist(path + "jobs/")
	if err != nil {
		return &rPipeline{}, err
	}

	repo := rPipeline{path}
	return &repo, nil
}

type rPipeline struct {
	basePath string
}
