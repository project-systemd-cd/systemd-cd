package pipeline

import (
	"systemd-cd/domain/logger"
	"systemd-cd/domain/unix"
)

func (p pipeline) test() (err error) {
	logger.Logger().Debug("START - Execute pipeline test command")
	defer func() {
		if err == nil {
			logger.Logger().Debug("END   - Execute pipeline test command")
		} else {
			logger.Logger().Error("FAILED - Execute pipeline test command")
			logger.Logger().Error(err)
		}
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
