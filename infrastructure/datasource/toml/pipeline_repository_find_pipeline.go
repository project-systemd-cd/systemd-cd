package toml

import (
	"bytes"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/toml"
	"systemd-cd/domain/unix"
)

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
