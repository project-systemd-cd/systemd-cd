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
	d := int64(time.Now().Unix() - int64(t))
	j.Duration = &d
	j.Logs = logs
	err2 := repo.SaveJob(j.Job)
	if err2 != nil {
		err = err2
		return err
	}

	return err
}

func (j jobInstance) Cancel(repo IRepository) (err error) {
	// Update job state
	j.Status = StatusJobCanceled
	err = repo.SaveJob(j.Job)
	return err
}
