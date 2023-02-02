package pipeline

import (
	"systemd-cd/domain/logger"
	"systemd-cd/domain/unix"
)

func (p pipeline) build() (err error) {
	logger.Logger().Info("-----------------------------------------------------------")
	logger.Logger().Info("START - Execute pipeline build command")
	logger.Logger().Infof("* pipeline.Name = %v", p.ManifestMerged.Name)
	logger.Logger().Tracef("* pipeline = %+v", p)
	logger.Logger().Info("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Info("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Info("END   - Execute pipeline build command")
		} else {
			logger.Logger().Error("FAILED - Execute pipeline build command")
			logger.Logger().Error(err)
		}
		logger.Logger().Info("-----------------------------------------------------------")
	}()

	if p.ManifestMerged.BuildCommands != nil {
		for _, cmd := range *p.ManifestMerged.BuildCommands {
			logger.Logger().Infof("Execute command \"%v\" (workingDir: \"%v\")", cmd, p.RepositoryLocal.Path)
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
