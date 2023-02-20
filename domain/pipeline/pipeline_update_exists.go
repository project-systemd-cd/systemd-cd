package pipeline

import (
	errorss "errors"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"
)

func (p *pipeline) updateExists() (exists bool, targetCommitId *string, targetTagName *string, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Check update")
	logger.Logger().Debugf("* pipeline.Name = %v", p.ManifestLocal.Name)
	logger.Logger().Tracef("* pipeline = %+v", *p)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debugf("> exists = %v", exists)
			if targetCommitId == nil {
				logger.Logger().Debug("> targetCommitId = nil")
			} else {
				logger.Logger().Debugf("> targetCommitId = %v", *targetCommitId)
			}
			logger.Logger().Debug("END   - Check update")
		} else {
			logger.Logger().Error("FAILED - Check update")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	err = p.RepositoryLocal.Fetch()
	if err != nil {
		return
	}

	if p.ManifestLocal.GitTagRegex != nil {
		var hash string
		var name string
		hash, name, err = p.RepositoryLocal.FindHashByTagRegex(*p.ManifestLocal.GitTagRegex)
		if err != nil {
			var ErrNotFound *errors.ErrNotFound
			if !errorss.As(err, &ErrNotFound) {
				return
			}
		} else {
			if hash != p.RepositoryLocal.RefCommitId {
				exists = true
				targetCommitId = &hash
				targetTagName = &name
			}
		}
	} else {
		var latest bool
		latest, err = p.RepositoryLocal.HeadIsLatesetOfBranch(p.ManifestLocal.GitTargetBranch)
		if err != nil {
			return
		}
		if !latest {
			exists = true
		}
	}

	return
}
