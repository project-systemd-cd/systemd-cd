package toml

import (
	"bytes"
	"strings"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/toml"
	"systemd-cd/domain/unix"
	"time"
)

// FindJobs implements pipeline.IRepository
func (r *rPipeline) FindJobs(pipelineName string, query pipeline.QueryParamJob) ([][]pipeline.Job, error) {
	lsOption := unix.LsOption{ReverceOrder: true, SortByDescendingTime: true, DirTrailiingSlash: true}
	if query.Asc {
		lsOption.ReverceOrder = false
	}
	s, err := unix.Ls(unix.ExecuteOption{}, lsOption, r.basePath+"jobs/")
	if err != nil {
		return nil, err
	}

	jobs := [][]pipeline.Job{}
	jobs2 := []pipeline.Job{}
	for _, v := range s {
		if strings.HasSuffix(v, "_"+pipelineName+".toml") {
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
			timestamp := time.Unix(int64(j.Timestamp), 0)

			if len(jobs2) != 0 && jobs2[0].GroupId != j.GroupId {
				jobs = append(jobs, jobs2)
				jobs2 = []pipeline.Job{}
			}
			if query.From == nil && query.To == nil {
				jobs2 = append(jobs2, j)
			} else if query.From != nil && !query.From.Before(timestamp) {
				jobs2 = append(jobs2, j)
			} else if query.To != nil && !query.To.After(timestamp) {
				jobs2 = append(jobs2, j)
			}
		}
	}
	if len(jobs2) != 0 {
		jobs = append(jobs, jobs2)
	}

	return jobs, nil
}
