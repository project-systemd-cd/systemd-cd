package toml

import (
	"bytes"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/toml"
	"systemd-cd/domain/unix"
)

// SavePipeline implements pipeline.IRepository
func (r *rPipeline) SavePipeline(m pipeline.PipelineMetadata) error {
	// Encode to toml format
	b := &bytes.Buffer{}
	err := toml.Encode(b, m, toml.EncodeOption{Indent: new(string)})
	if err != nil {
		return err
	}

	// Write to file
	err = unix.WriteFile(r.basePath+m.Name+".toml", b.Bytes())
	if err != nil {
		return err
	}

	// If git remote url changed, delete pipeline job data
	j, err := r.FindJobs(m.Name, pipeline.QueryParamJob{})
	if err != nil {
		return err
	}
	if len(j) != 0 && len(j[0]) != 0 && j[0][0].GitRemoteUrl != m.ManifestLocal.GitRemoteUrl {
		err = unix.Rm(
			unix.ExecuteOption{WantExitCodes: []int{1}},
			unix.RmOption{},
			r.basePath+"jobs/*_"+m.Name+".toml",
		)
		if err != nil {
			return err
		}
	}

	return nil
}
