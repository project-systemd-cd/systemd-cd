package pipeline

import (
	errorss "errors"
	"fmt"
	"strings"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/systemd"
	"systemd-cd/domain/unix"
	"time"
)

func (p pipeline) newJobInstall(groupId string, tag *string) (job *jobInstance, err error) {
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

	var commitAuthor string
	commitAuthor, err = p.GetCommitAuthor()
	if err != nil {
		return nil, err
	}
	var commitMsg string
	commitMsg, err = p.GetCommitMessage()
	if err != nil {
		return nil, err
	}
	job = &jobInstance{
		Job: Job{
			GroupId:       groupId,
			Id:            UUID(),
			PipelineName:  p.ManifestMerged.Name,
			GitRemoteUrl:  p.ManifestLocal.GitRemoteUrl,
			Branch:        p.ManifestMerged.GitTargetBranch,
			Tag:           tag,
			CommitId:      p.GetCommitRef(),
			CommitAuthor:  commitAuthor,
			CommitMessage: commitMsg,
			Type:          JobTypeInstall,
			Status:        StatusJobPending,
		},
	}

	job.f = func() (logs []jobLog, err2 error) {
		logger.Logger().Debug("-----------------------------------------------------------")
		logger.Logger().Debug("START - Install pipeline files")
		logger.Logger().Debugf("* pipeline.Name = %v", p.ManifestMerged.Name)
		logger.Logger().Tracef("* pipeline = %+v", p)
		logger.Logger().Debug("-----------------------------------------------------------")
		defer func() {
			logger.Logger().Debug("-----------------------------------------------------------")
			if err2 == nil {
				logger.Logger().Debug("END   - Install pipeline files")
			} else {
				logger.Logger().Error("FAILED - Install pipeline files")
				logger.Logger().Error(err2)
			}
			logger.Logger().Debug("-----------------------------------------------------------")
		}()

		logger.Logger().Infof("Install files of pipeline \"%s\" (version: \"%s\")", p.ManifestMerged.Name, p.GetCommitRef())

		if p.ManifestMerged.Binaries != nil && len(*p.ManifestMerged.Binaries) != 0 {
			logger.Logger().Debug("Install binary files")
			pathBinDir := p.service.PathBinDir + p.ManifestMerged.Name + "/"
			err2 = unix.MkdirIfNotExist(pathBinDir)
			if err2 != nil {
				return logs, err2
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
					return logs, err2
				}

				logs = append(logs, log)
			}
		}

		if p.ManifestMerged.SystemdServiceOptions != nil && len(p.ManifestMerged.SystemdServiceOptions) != 0 {
			for _, service := range p.ManifestMerged.SystemdServiceOptions {
				logger.Logger().Debugf("Install files for \"%s\" as systemd unit", service.Name)
				if service.Etc != nil {
					logger.Logger().Debug(" Install etc files")
					pathEtcDir := p.service.PathEtcDir + service.Name + "/"
					err2 = unix.MkdirIfNotExist(pathEtcDir)
					if err2 != nil {
						return logs, err2
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
								return logs, err2
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
								return logs, err2
							}
						}

						logs = append(logs, log)
					}
				}

				if service.Opt != nil {
					logger.Logger().Debug(" Install opt files")
					pathOptDir := p.service.PathOptDir + service.Name + "/"
					err2 = unix.MkdirIfNotExist(pathOptDir)
					if err2 != nil {
						return logs, err2
					}

					// Copy opt files
					for _, src := range service.Opt {
						dest := pathOptDir + src
						cpOption := unix.CpOption{Recursive: true, Force: true}
						if strings.Contains(strings.TrimSuffix(src, "/"), "/") {
							dest = pathOptDir
							cpOption.Parents = true
						}

						log := jobLog{Commmand: "cp -R --parents -f" + src + " " + dest}

						// Copy opt file
						err2 = unix.Cp(
							unix.ExecuteOption{WorkingDirectory: (*string)(&p.RepositoryLocal.Path)},
							cpOption,
							src,
							dest,
						)
						if err2 != nil {
							log.Output = err2.Error()
							logs = append(logs, log)
							return logs, err2
						}

						logs = append(logs, log)
					}
				}

				{
					u := service.generateSystemdServiceUnit(&p)

					var s systemd.IUnitService
					s, err2 := p.service.Systemd.GetService(u.Name)
					notInstalled := false
					var ErrNotFound *errors.ErrNotFound
					if err == nil {
						logger.Logger().Info(" Upgrate systemd unit")
					} else if errorss.As(err, &ErrNotFound) {
						logger.Logger().Info(" Install systemd unit")
						notInstalled = true
					} else {
						return logs, err2
					}

					s, err2 = p.service.Systemd.NewService(u.Name, u.UnitFile, u.Env)
					var b []byte
					b, _ = systemd.MarshalUnitFile(u.UnitFile)
					log := jobLog{Commmand: fmt.Sprintf("cat << EOF > %s\n%s\nEOF", s.GetUnitFilePath(), string(b))}
					if err2 != nil {
						log.Output = err2.Error()
						logs = append(logs, log)
						return logs, err2
					}
					logs = append(logs, log)

					logger.Logger().Info(" Run systemd unit")
					if notInstalled {
						log := jobLog{Commmand: fmt.Sprintf("systemctl enable %s.service", u.Name)}
						err2 = s.Enable(true)
						if err2 != nil {
							log.Output = err2.Error()
							logs = append(logs, log)
							return logs, err2
						}
						logs = append(logs, log)
					} else {
						log := jobLog{Commmand: fmt.Sprintf("systemctl restart %s.service", u.Name)}
						err2 = s.Restart()
						if err2 != nil {
							log.Output = err2.Error()
							logs = append(logs, log)
							return logs, err2
						}
						logs = append(logs, log)
					}

					time.Sleep(time.Second)

					// Get status of systemd service
					var status systemd.Status
					status, err2 = s.GetStatus()
					if err2 != nil {
						return logs, err2
					}
					if status == systemd.StatusFailed {
						return logs, errorss.New("failed to run systemd unit")
					}
				}
			}
		}

		return logs, err2
	}

	err = p.service.repo.SaveJob(job.Job)

	return job, err
}
