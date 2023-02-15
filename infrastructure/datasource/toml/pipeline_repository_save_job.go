package toml

import (
	"bytes"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/toml"
	"systemd-cd/domain/unix"
)

// SaveJob implements pipeline.IRepository
func (r *rPipeline) SaveJob(job pipeline.Job) error {
	b := &bytes.Buffer{}

	// Encode to toml format
	err := toml.Encode(b, job, toml.EncodeOption{Indent: new(string)})
	if err != nil {
		return err
	}

	// Write to file
	err = unix.WriteFile(r.basePath+"jobs/"+job.GroupId+"_"+job.Id+"_"+job.PipelineName+".toml", b.Bytes())
	if err != nil {
		return err
	}

	return nil
}
