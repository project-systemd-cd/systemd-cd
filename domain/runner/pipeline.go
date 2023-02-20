package runner

import "systemd-cd/domain/pipeline"

type Pipeline struct {
	pipeline.IPipeline
	autoSyncEnabled bool
}

func (p Pipeline) AutoSyncEnabled() bool {
	return p.autoSyncEnabled
}
