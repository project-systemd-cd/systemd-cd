package toml

import (
	"bytes"
	"strings"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/toml"
	"systemd-cd/domain/unix"
	"time"
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

// FindJob implements pipeline.IRepository
func (r *rPipeline) FindJob(groupId string) ([]pipeline.Job, error) {
	wd := r.basePath + "jobs/"
	s, err := unix.Ls(
		unix.ExecuteOption{WorkingDirectory: &wd},
		unix.LsOption{ReverceOrder: true, SortByDescendingTime: true, DirTrailiingSlash: true},
		groupId+"_*.toml",
	)
	if err != nil {
		if !strings.Contains(err.Error(), "No such file or directory") {
			return nil, err
		}
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
