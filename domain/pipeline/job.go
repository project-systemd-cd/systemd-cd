package pipeline

import (
	"time"
)

type Job struct {
	PipeineId string

	Id string

	PipelineName string
	CommitId     string
	Type         jobType

	Status statusJob

	Timestamp unixTime
	Duration  *int64
	Logs      []jobLog
}

type jobType string

const (
	JobTypeTest    jobType = "test"
	JobTypeBuild   jobType = "build"
	JobTypeInstall jobType = "install"
)

type statusJob string

const (
	StatusJobPending    statusJob = "pending"
	StatusJobDone       statusJob = "done"
	StatusJobInProgress statusJob = "in progress"
	StatusJobFailed     statusJob = "failed"
	StatusJobCanceled   statusJob = "canceled"
)

type unixTime int64

type jobLog struct {
	Commmand string
	Output   string
}

type UpdateParamJob struct {
	Status   *statusJob
	Duration *time.Duration
	Stdout   *string
}

type QueryParamJob struct {
	From *time.Time
	To   *time.Time
	Asc  bool
}
