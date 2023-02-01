package pipeline

import (
	"systemd-cd/domain/logger"
	"systemd-cd/domain/unix"
)

func (p pipeline) build() (err error) {
	logger.Logger().Debug("START - Execute pipeline build command")
	logger.Logger().Debugf("< pipeline.Name = %v", p.ManifestMerged.Name)
	defer func() {
		if err == nil {
			logger.Logger().Debug("END   - Execute pipeline build command")
		} else {
			logger.Logger().Error("FAILED - Execute pipeline build command")
			logger.Logger().Error(err)
		}
	}()

	if p.ManifestMerged.BuildCommands != nil {
		for _, cmd := range *p.ManifestMerged.BuildCommands {
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
