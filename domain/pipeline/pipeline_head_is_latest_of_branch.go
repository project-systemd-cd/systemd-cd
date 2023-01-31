package pipeline

import "systemd-cd/domain/logger"

func (p pipeline) HeadIsLatesetOfBranch(branch string) (bool, error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: p}}))

	// Check update
	isLatest, err := p.RepositoryLocal.HeadIsLatesetOfBranch(branch)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return false, nil
	}

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "isLatest", Value: isLatest}}))
	return isLatest, nil
}
