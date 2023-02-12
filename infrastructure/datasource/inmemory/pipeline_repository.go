package inmemory

import (
	"systemd-cd/domain/errors"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/runner"
)

func NewRepositoryPipelineInmemory() runner.IRepositoryInmemory {
	return &rPipeline{new([]pipeline.IPipeline)}
}

type rPipeline struct {
	pipelines *[]pipeline.IPipeline
}

// AddPipeline implements runner.IRepositoryInmemory
func (r *rPipeline) AddPipeline(new pipeline.IPipeline) (pipeline.IPipeline, error) {
	*r.pipelines = append(*r.pipelines, new)
	return new, nil
}

// FindPipeline implements runner.IRepositoryInmemory
func (r *rPipeline) FindPipeline(name string) (pipeline.IPipeline, error) {
	var err error = &errors.ErrNotFound{Object: "pipeline", IdName: "name", Id: name}
	for _, p := range *r.pipelines {
		if p.GetName() == name {
			return p, nil
		}
	}
	return nil, err
}

// FindPipelines implements runner.IRepositoryInmemory
func (r *rPipeline) FindPipelines() ([]pipeline.IPipeline, error) {
	return *r.pipelines, nil
}
