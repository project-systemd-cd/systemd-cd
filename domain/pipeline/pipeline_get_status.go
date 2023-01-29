package pipeline

import "systemd-cd/domain/logger"

func (p pipeline) GetStatus() Status {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: p}}))
	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "Status", Value: p.Status}}))
	return p.Status
}
