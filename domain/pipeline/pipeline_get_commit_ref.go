package pipeline

import "systemd-cd/domain/logger"

func (p pipeline) GetCommitRef() string {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: p}}))
	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "RefCommitId", Value: p.RepositoryLocal.RefCommitId}}))
	return p.RepositoryLocal.RefCommitId
}
