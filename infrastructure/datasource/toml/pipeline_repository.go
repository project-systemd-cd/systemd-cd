package toml

import (
	"bytes"
	"strings"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/toml"
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

	repo := rPipeline{path}
	return &repo, nil
}

type rPipeline struct {
	basePath string
}

// FindPipelineByName implements pipeline.IRepository
func (r *rPipeline) FindPipelineByName(name string) (pipeline.PipelineMetadata, error) {
	s, err := unix.Ls(
		unix.ExecuteOption{},
		unix.LsOption{DirTrailiingSlash: true},
		r.basePath,
	)
	if err != nil {
		return pipeline.PipelineMetadata{}, err
	}

	for _, v := range s {
		if v == name+".toml" {
			// Read file
			b := &bytes.Buffer{}
			err = unix.ReadFile(r.basePath+v, b)
			if err != nil {
				return pipeline.PipelineMetadata{}, err
			}

			// Unmarshal toml
			m := pipeline.PipelineMetadata{}
			err = toml.Decode(b, &m)
			if err != nil {
				return pipeline.PipelineMetadata{}, err
			}

			return m, nil
		}
	}

	err = &errors.ErrNotFound{Object: "Pipeline", IdName: "name", Id: name}
	return pipeline.PipelineMetadata{}, err
}

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

// SavePipeline implements pipeline.IRepository
func (r *rPipeline) SavePipeline(m pipeline.PipelineMetadata) error {
	b := &bytes.Buffer{}

	// Encode to toml format
	err := toml.Encode(b, m, toml.EncodeOption{Indent: new(string)})
	if err != nil {
		return err
	}

	// Write to file
	err = unix.WriteFile(r.basePath+m.Name+".toml", b.Bytes())
	if err != nil {
		return err
	}

	return nil
}
