package git

import "systemd-cd/domain/logger"

func (r *RepositoryLocal) FindHashByTagRegex(regex string) (hash string, err error) {
	logger.Logger().Debug("START - Find git commit hash by regex of git tag")
	defer func() {
		if err == nil {
			logger.Logger().Debug("END   - Find git commit hash by regex of git tag")
		} else {
			logger.Logger().Error("FAILED - Find git commit hash by regex of git tag")
			logger.Logger().Error(err)
		}
	}()

	hash, err = r.git.command.FindHashByTagRegex(r.Path, regex)
	return hash, err
}
