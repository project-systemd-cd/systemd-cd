package pipeline

import (
	"time"
)

type Job struct {
	GroupId string `json:"group_id"`

	Id string `json:"id"`

	PipelineName string `json:"pipeline_name"`
	GitRemoteUrl string `json:"git_remote_url"`

	Branch        string  `json:"branch"`
	Tag           *string `json:"tag,omitempty"`
	CommitId      string  `json:"commit_id"`
	CommitAuthor  string  `json:"commit_author"`
	CommitMessage string  `json:"commit_message"`
	Type          jobType `json:"type"`

	Status statusJob `json:"status"`

	Timestamp *unixTime `json:"timestamp,omitempty"`
	Duration  *int64    `json:"duration"`
	Logs      []jobLog  `json:"logs,omitempty"`
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
	Commmand string `json:"command"`
	Output   string `json:"output"`
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
