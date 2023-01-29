package pipeline

import "systemd-cd/domain/logger"

func (p *pipeline) CheckUpdate() (updateExists bool, err error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: p}}))

	// Check update
	updateExists, err = p.RepositoryLocal.DiffExists(true)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return
	}
	if updateExists {
		p.Status = StatusOutOfSync
	} else {
		p.Status = StatusSynced
	}

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "updateExists", Value: updateExists}}))
	return updateExists, nil
}
