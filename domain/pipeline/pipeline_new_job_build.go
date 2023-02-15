package pipeline

import (
	"bytes"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/unix"
)

func (p pipeline) newJobBuild(groupId string) (job *jobInstance, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Register job for build")
	logger.Logger().Debugf("* pipeline.Name = %v", p.ManifestMerged.Name)
	logger.Logger().Tracef("* pipeline = %+v", p)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			if job != nil {
				logger.Logger().Debugf("> job.Id = %s", job.Id)
				logger.Logger().Tracef("> job = %+v", job)
			} else {
				logger.Logger().Debug("> job = nil")
			}
			logger.Logger().Debug("END   - Register job for build")
		} else {
			logger.Logger().Error("FAILED - Register job for build")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	if p.ManifestMerged.BuildCommands != nil {
		job = &jobInstance{
			Job: Job{
				GroupId:      groupId,
				Id:           UUID(),
				PipelineName: p.ManifestMerged.Name,
				CommitId:     p.GetCommitRef(),
				Type:         JobTypeBuild,
				Status:       StatusJobPending,
			},
		}

		job.f = func() (logs []jobLog, err2 error) {
			logger.Logger().Info("-----------------------------------------------------------")
			logger.Logger().Info("START - Execute pipeline build command")
			logger.Logger().Infof("* pipeline.Name = %v", p.ManifestMerged.Name)
			logger.Logger().Infof("* job.Id = %v", job.Id)
			logger.Logger().Tracef("* pipeline = %+v", p)
			logger.Logger().Info("-----------------------------------------------------------")
			defer func() {
				logger.Logger().Info("-----------------------------------------------------------")
				if err2 == nil {
					logger.Logger().Info("END   - Execute pipeline build command")
				} else {
					logger.Logger().Error("FAILED - Execute pipeline build command")
					logger.Logger().Error(err2)
				}
				logger.Logger().Info("-----------------------------------------------------------")
			}()

			for _, cmd := range *p.ManifestMerged.BuildCommands {
				logger.Logger().Infof("Execute command \"%v\" (workingDir: \"%v\")", cmd, p.RepositoryLocal.Path)
				log := jobLog{Commmand: cmd}

				var stdout bytes.Buffer
				_, stdout, _, err2 = unix.Execute(
					unix.ExecuteOption{
						WorkingDirectory: (*string)(&p.RepositoryLocal.Path),
					},
					"/usr/bin/bash", "-c", "\""+cmd+"\"",
				)
				if err2 != nil {
					log.Output = err2.Error()
					logs = append(logs, log)

					return logs, err2
				}

				log.Output = stdout.String()
				logs = append(logs, log)
			}

			return logs, err2
		}

		err = p.service.repo.SaveJob(job.Job)
	}

	return job, err
}
