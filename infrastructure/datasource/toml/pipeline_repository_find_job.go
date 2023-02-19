package toml

import (
	"bytes"
	"strings"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/toml"
	"systemd-cd/domain/unix"
)

// FindJob implements pipeline.IRepository
func (r *rPipeline) FindJob(groupId string) ([]pipeline.Job, error) {
	wd := r.basePath + "jobs/"
	s, err := unix.Ls(
		unix.ExecuteOption{WorkingDirectory: &wd},
		unix.LsOption{ReverceOrder: false, SortByDescendingTime: true, DirTrailiingSlash: true},
		groupId+"_*.toml",
	)
	if err != nil {
		if !strings.Contains(err.Error(), "No such file or directory") {
			return nil, err
		}
		err = &errors.ErrNotFound{
			Object: "PipelineJobGroup",
			IdName: "groupId",
			Id:     groupId,
		}
		return nil, err
	}

	jobs := []pipeline.Job{}
	for _, v := range s {
		// Read file
		b := &bytes.Buffer{}
		err = unix.ReadFile(r.basePath+"jobs/"+v, b)
		if err != nil {
			return nil, err
		}

		// Unmarshal toml
		j := pipeline.Job{}
		err = toml.Decode(b, &j)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, j)
	}

	return jobs, nil
}
