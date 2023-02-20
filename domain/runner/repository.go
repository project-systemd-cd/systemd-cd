package runner

import (
	"systemd-cd/domain/errors"
	"systemd-cd/domain/pipeline"
)

type iRepositoryInmemory interface {
	AddPipeline(pipeline.IPipeline) (pipeline.IPipeline, error)

	FindPipeline(name string) (pipeline.IPipeline, error)
	FindPipelines() ([]pipeline.IPipeline, error)
}

func inmemoryRepository() iRepositoryInmemory {
	return &rPipeline{new([]pipeline.IPipeline)}
}

type rPipeline struct {
	pipelines *[]pipeline.IPipeline
}

func (r *rPipeline) AddPipeline(new pipeline.IPipeline) (pipeline.IPipeline, error) {
	*r.pipelines = append(*r.pipelines, new)
	return new, nil
}

func (r *rPipeline) FindPipeline(name string) (pipeline.IPipeline, error) {
	var err error = &errors.ErrNotFound{Object: "pipeline", IdName: "name", Id: name}
	for _, p := range *r.pipelines {
		if p.GetName() == name {
			return p, nil
		}
	}
	return nil, err
}

func (r *rPipeline) FindPipelines() ([]pipeline.IPipeline, error) {
	return *r.pipelines, nil
}
