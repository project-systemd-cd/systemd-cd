package toml

import (
	"bytes"
	"strings"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/toml"
	"systemd-cd/domain/unix"
)

func NewRepositoryPipeline(path string) (pipeline.IRepository, error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "path", Value: path}}))

	err := unix.MkdirIfNotExist(path)
	if err != nil {
		return &rPipeline{}, err
	}
	if !strings.HasSuffix(path, "/") {
		// Add trailing slash
		path += "/"
	}

	repo := rPipeline{path}
	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Value: repo}}))
	return &repo, nil
}

type rPipeline struct {
	basePath string
}

// FindPipeline implements pipeline.IRepository
func (r *rPipeline) FindPipeline(name string) (pipeline.PipelineMetadata, error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "name", Value: name}}))

	s, err := unix.Ls(
		unix.ExecuteOption{},
		unix.LsOption{DirTrailiingSlash: true},
		r.basePath,
	)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return pipeline.PipelineMetadata{}, err
	}

	for _, v := range s {
		if v == name+".toml" {
			// Read file
			b := &bytes.Buffer{}
			err = unix.ReadFile(r.basePath+v, b)
			if err != nil {
				logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
				return pipeline.PipelineMetadata{}, err
			}

			// Unmarshal toml
			m := pipeline.PipelineMetadata{}
			err = toml.Decode(b, &m)
			if err != nil {
				logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
				return pipeline.PipelineMetadata{}, err
			}

			logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Value: m}}))
			return m, nil
		}
	}

	err = &errors.ErrNotFound{Object: "Pipeline", IdName: "name", Id: name}
	logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
	return pipeline.PipelineMetadata{}, err
}

// FindPipelines implements pipeline.IRepository
func (r *rPipeline) FindPipelines() (pipelines []pipeline.PipelineMetadata, err error) {
	logger.Logger().Trace("Called")

	s, err := unix.Ls(
		unix.ExecuteOption{},
		unix.LsOption{DirTrailiingSlash: true},
		r.basePath,
	)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
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
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return nil, err
		}

		// Decode toml
		m := pipeline.PipelineMetadata{}
		err = toml.Decode(b, &m)
		if err != nil {
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return nil, err
		}

		pipelines = append(pipelines, m)
	}

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Value: pipelines}}))
	return pipelines, nil
}

// SavePipeline implements pipeline.IRepository
func (r *rPipeline) SavePipeline(m pipeline.PipelineMetadata) error {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: m}}))

	b := &bytes.Buffer{}

	// Encode to toml format
	err := toml.Encode(b, m, toml.EncodeOption{Indent: new(string)})
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	// Write to file
	err = unix.WriteFile(r.basePath+m.Name+".toml", b.Bytes())
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}
