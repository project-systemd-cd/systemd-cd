package pipeline

import (
	"systemd-cd/domain/model/logger"
	"systemd-cd/domain/model/unix"
)

func (p pipeline) test() error {
	logger.Logger().Tracef("Called:\n\tpipeline: %+v", p)

	if p.ManifestMerged.TestCommand != nil {
		_, _, _, err := unix.Execute(
			unix.ExecuteOption{
				WorkingDirectory: (*string)(&p.RepositoryLocal.Path),
			},
			"/usr/bin/bash", "-c", "\""+*p.ManifestMerged.TestCommand+"\"",
		)
		if err != nil {
			logger.Logger().Errorf("Error:\n\terr: %v", err)
			return err
		}
	}

	logger.Logger().Trace("Finished")
	return nil
}
