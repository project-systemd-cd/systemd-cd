package git

import (
	errorss "errors"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"
)

func (r *RepositoryLocal) FindHashByTagRegex(regex string) (hash string, err error) {
	logger.Logger().Debug("START - Find git commit hash by regex of git tag")
	logger.Logger().Debugf("< regex = %v", regex)
	defer func() {
		var ErrNotFound *errors.ErrNotFound
		if err == nil || errorss.As(err, &ErrNotFound) {
			logger.Logger().Debugf("> commitId = %v", hash)
			logger.Logger().Debug("END   - Find git commit hash by regex of git tag")
		} else {
			logger.Logger().Error("FAILED - Find git commit hash by regex of git tag")
			logger.Logger().Error(err)
		}
	}()

	hash, err = r.git.command.FindHashByTagRegex(r.Path, regex)
	return hash, err
}
