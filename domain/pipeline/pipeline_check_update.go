package pipeline

import "systemd-cd/domain/logger"

func (p pipeline) GetUpdateExistence() (updateExists bool, err error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: p}}))

	// Check update
	updateExists, err = p.RepositoryLocal.DiffExists(true)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return
	}

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "updateExists", Value: updateExists}}))
	return updateExists, nil
}
