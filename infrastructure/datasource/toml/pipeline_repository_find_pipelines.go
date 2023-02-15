package toml

import (
	"bytes"
	"strings"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/toml"
	"systemd-cd/domain/unix"
)

// FindPipelines implements pipeline.IRepository
func (r *rPipeline) FindPipelines() (pipelines []pipeline.PipelineMetadata, err error) {
	s, err := unix.Ls(
		unix.ExecuteOption{},
		unix.LsOption{DirTrailiingSlash: true},
		r.basePath,
	)
	if err != nil {
		return nil, err
	}

	for _, v := range s {
		if !strings.HasSuffix(v, ".toml") {
			// if not toml file, skip
			continue
		}

		// Read file
		b := &bytes.Buffer{}
		err = unix.ReadFile(r.basePath+v, b)
		if err != nil {
			return nil, err
		}

		// Decode toml
		m := pipeline.PipelineMetadata{}
		err = toml.Decode(b, &m)
		if err != nil {
			return nil, err
		}

		pipelines = append(pipelines, m)
	}

	return pipelines, nil
}
