package runner

import (
	"systemd-cd/domain/errors"
)

type iRepositoryInmemory interface {
	AddPipeline(Pipeline) (Pipeline, error)

	FindPipeline(name string) (Pipeline, error)
	FindPipelines() ([]Pipeline, error)

	RemovePipeline(name string) error
}

func inmemoryRepository() iRepositoryInmemory {
	return &rPipeline{new([]Pipeline)}
}

type rPipeline struct {
	pipelines *[]Pipeline
}

func (r *rPipeline) AddPipeline(new Pipeline) (Pipeline, error) {
	*r.pipelines = append(*r.pipelines, new)
	return new, nil
}

func (r *rPipeline) FindPipeline(name string) (Pipeline, error) {
	var err error = &errors.ErrNotFound{Object: "pipeline", IdName: "name", Id: name}
	for _, p := range *r.pipelines {
		if p.GetName() == name {
			return p, nil
		}
	}
	return Pipeline{}, err
}

func (r *rPipeline) FindPipelines() ([]Pipeline, error) {
	return *r.pipelines, nil
}

func (r *rPipeline) RemovePipeline(name string) error {
	var err error = &errors.ErrNotFound{Object: "pipeline", IdName: "name", Id: name}
	newPipelines := &[]Pipeline{}
	for _, p := range *r.pipelines {
		if p.GetName() == name {
			err = nil
		} else {
			(*newPipelines) = append((*newPipelines), p)
		}
	}
	if err == nil {
		r.pipelines = newPipelines
	}
	return err
}
