package pipeline

import (
	"fmt"
	"strings"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/systemd"
	"systemd-cd/domain/unix"
)

func (p pipeline) newJobInstall(groupId string) (job *jobInstance, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Register job for install")
	logger.Logger().Debugf("* pipeline.Name = %v", p.ManifestMerged.Name)
	logger.Logger().Tracef("* pipeline = %+v", p)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debugf("> job.Id = %s", job.Id)
			logger.Logger().Tracef("> job = %+v", job)
			logger.Logger().Debug("END   - Register job for install")
		} else {
			logger.Logger().Error("FAILED - Register job for install")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	job = &jobInstance{
		Job: Job{
			GroupId:           groupId,
			Id:                UUID(),
			PipelineName:      p.ManifestMerged.Name,
			GitTargetBranch:   p.ManifestMerged.GitTargetBranch,
			GitTargetTagRegex: p.ManifestMerged.GitTagRegex,
			CommitId:          p.GetCommitRef(),
			Type:              JobTypeInstall,
			Status:            StatusJobPending,
		},
	}

	job.f = func() (logs []jobLog, err2 error) {
		logger.Logger().Info("-----------------------------------------------------------")
		logger.Logger().Info("START - Install pipeline files")
		logger.Logger().Infof("* pipeline.Name = %v", p.ManifestMerged.Name)
		logger.Logger().Tracef("* pipeline = %+v", p)
		logger.Logger().Info("-----------------------------------------------------------")
		defer func() {
			logger.Logger().Info("-----------------------------------------------------------")
			if err2 == nil {
				logger.Logger().Info("END   - Install pipeline files")
			} else {
				logger.Logger().Error("FAILED - Install pipeline files")
				logger.Logger().Error(err2)
			}
			logger.Logger().Info("-----------------------------------------------------------")
		}()

		if p.ManifestMerged.Binaries != nil && len(*p.ManifestMerged.Binaries) != 0 {
			logger.Logger().Info("Install binary files")
			pathBinDir := p.service.PathBinDir + p.ManifestMerged.Name + "/"
			err2 = unix.MkdirIfNotExist(pathBinDir)
			if err2 != nil {
				return nil, err2
			}
			for _, binary := range *p.ManifestMerged.Binaries {
				src := string(p.RepositoryLocal.Path) + "/" + binary
				dest := pathBinDir + strings.TrimPrefix(binary, "./")

				log := jobLog{Commmand: "cp -f " + src + " " + dest}

				// Copy binary files
				err2 = unix.Cp(unix.ExecuteOption{}, unix.CpOption{Force: true}, src, dest)
				if err2 != nil {
					log.Output = err2.Error()
					logs = append(logs, log)
					return nil, err2
				}

				logs = append(logs, log)
			}
		}

		if p.ManifestMerged.SystemdServiceOptions != nil && len(p.ManifestMerged.SystemdServiceOptions) != 0 {
			for _, service := range p.ManifestMerged.SystemdServiceOptions {
				logger.Logger().Infof("Install files for \"%s\" as systemd unit", service.Name)
				if service.Etc != nil {
					logger.Logger().Info(" Install etc files")
					pathEtcDir := p.service.PathEtcDir + service.Name + "/"
					err2 = unix.MkdirIfNotExist(pathEtcDir)
					if err2 != nil {
						return nil, err2
					}

					// Copy or create etc files and add to cli options
					for _, etc := range service.Etc {
						src := strings.TrimPrefix(etc.Target, "./")
						dest := pathEtcDir + src

						log := jobLog{}

						if etc.Content != nil {
							log.Commmand = fmt.Sprintf("cat << EOF > %s\n%s\nEOF", dest, strings.ReplaceAll(*etc.Content, "'", "\\'"))

							// Create etc file
							err2 = unix.WriteFile(dest, []byte(*etc.Content))
							if err2 != nil {
								log.Output = err2.Error()
								logs = append(logs, log)
								return nil, err2
							}
						} else {
							log.Commmand = "cp -Rf " + src + " " + dest

							// Copy etc file
							err2 = unix.Cp(
								unix.ExecuteOption{WorkingDirectory: (*string)(&p.RepositoryLocal.Path)},
								unix.CpOption{Recursive: true, Force: true},
								src, dest,
							)
							if err2 != nil {
								log.Output = err2.Error()
								logs = append(logs, log)
								return nil, err2
							}
						}

						logs = append(logs, log)
					}
				}

				if service.Opt != nil {
					logger.Logger().Info(" Install opt files")
					pathOptDir := p.service.PathOptDir + service.Name + "/"
					err2 = unix.MkdirIfNotExist(pathOptDir)
					if err2 != nil {
						return nil, err2
					}

					// Copy opt files
					for _, src := range service.Opt {
						dest := pathOptDir + src

						log := jobLog{Commmand: "cp -RPf " + src + " " + dest}

						// Copy opt file
						err2 = unix.Cp(
							unix.ExecuteOption{WorkingDirectory: (*string)(&p.RepositoryLocal.Path)},
							unix.CpOption{Recursive: true, Parents: true, Force: true},
							src,
							dest,
						)
						if err2 != nil {
							log.Output = err2.Error()
							logs = append(logs, log)
							return nil, err2
						}

						logs = append(logs, log)
					}
				}

				{
					logger.Logger().Info(" Install systemd service unit file")

					u := service.generateSystemdServiceUnit(&p)

					var s systemd.UnitService
					s, err2 = p.service.Systemd.NewService(u.Name, u.UnitFile, u.Env)

					var b []byte
					b, _ = systemd.MarshalUnitFile(u.UnitFile)
					log := jobLog{Commmand: fmt.Sprintf("cat << EOF > %s\n%s\nEOF", s.Path, string(b))}

					if err2 != nil {
						log.Output = err2.Error()
						logs = append(logs, log)
						return nil, err2
					}

					logs = append(logs, log)
				}
			}
		}

		return logs, err2
	}

	err = p.service.repo.SaveJob(job.Job)

	return job, err
}
