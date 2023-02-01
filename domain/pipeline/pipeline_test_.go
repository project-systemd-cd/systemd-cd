package pipeline

import (
	"systemd-cd/domain/logger"
	"systemd-cd/domain/unix"
)

func (p pipeline) test() (err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Execute pipeline test command")
	logger.Logger().Debugf("* pipeline.Name = %v", p.ManifestMerged.Name)
	logger.Logger().Tracef("* pipeline = %+v", p)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debug("END   - Execute pipeline test command")
		} else {
			logger.Logger().Error("FAILED - Execute pipeline test command")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	if p.ManifestMerged.TestCommands != nil {
		for _, cmd := range *p.ManifestMerged.TestCommands {
			_, _, _, err = unix.Execute(
				unix.ExecuteOption{
					WorkingDirectory: (*string)(&p.RepositoryLocal.Path),
				},
				"/usr/bin/bash", "-c", "\""+cmd+"\"",
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
