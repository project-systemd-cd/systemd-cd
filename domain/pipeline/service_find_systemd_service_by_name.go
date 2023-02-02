package pipeline

import "systemd-cd/domain/systemd"

// FindSystemdServiceByName implements IPipelineService
func (s pipelineService) FindSystemdServiceByName(name string) (systemd.IService, error) {
	// TODO: unimplemented
	panic("unimplemented")
}
