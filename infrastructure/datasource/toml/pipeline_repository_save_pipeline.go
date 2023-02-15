package toml

import (
	"bytes"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/toml"
	"systemd-cd/domain/unix"
)

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
