package runner

import "systemd-cd/domain/pipeline"

type IRepositoryInmemory interface {
	AddPipeline(pipeline.IPipeline) (pipeline.IPipeline, error)

	FindPipeline(name string) (pipeline.IPipeline, error)
	FindPipelines() ([]pipeline.IPipeline, error)
}
