package pipeline

import "time"

type jobInstance struct {
	Job
	f func() (logs []jobLog, err error)
}

func (j jobInstance) Run(repo IRepository) (err error) {
	// Update job state
	startAt := time.Now()
	t := unixTime(startAt.Unix())
	j.Timestamp = &t
	j.Status = StatusJobInProgress
	err = repo.SaveJob(j.Job)
	if err != nil {
		return err
	}

	// Run job
	logs, err := j.f()

	// Update job state
	j.Status = StatusJobDone
	if err != nil {
		j.Status = StatusJobFailed
	}
	d := int64(time.Since(startAt))
	j.Duration = &d
	j.Logs = logs
	err = repo.SaveJob(j.Job)
	if err != nil {
		return err
	}

	return err
}
