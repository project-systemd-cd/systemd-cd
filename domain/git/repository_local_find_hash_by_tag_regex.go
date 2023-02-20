package git

import (
	errorss "errors"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"
)

func (r *RepositoryLocal) FindHashByTagRegex(regex string) (hash string, name string, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Find git commit hash by regex of git tag")
	logger.Logger().Debugf("< regex = %v", regex)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		var ErrNotFound *errors.ErrNotFound
		if err == nil || errorss.As(err, &ErrNotFound) {
			logger.Logger().Debugf("> commitId = %v", hash)
			logger.Logger().Debug("END   - Find git commit hash by regex of git tag")
		} else {
			logger.Logger().Error("FAILED - Find git commit hash by regex of git tag")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	hash, name, err = r.git.command.FindHashByTagRegex(r.Path, regex)
	return hash, name, err
}
