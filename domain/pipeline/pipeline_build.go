package pipeline

import (
	"systemd-cd/domain/logger"
	"systemd-cd/domain/unix"
)

func (p pipeline) build() error {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: p}}))

	if p.ManifestMerged.BuildCommands != nil {
		for _, cmd := range *p.ManifestMerged.BuildCommands {
			_, _, _, err := unix.Execute(
				unix.ExecuteOption{
					WorkingDirectory: (*string)(&p.RepositoryLocal.Path),
				},
				"/usr/bin/bash", "-c", "\""+cmd+"\"",
			)
			if err != nil {
				logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
				return err
			}
		}
	}

	logger.Logger().Trace("Finished")
	return nil
}
