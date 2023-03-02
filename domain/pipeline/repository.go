package pipeline

import (
	"systemd-cd/domain/git"
)

type IRepository interface {
	SavePipeline(PipelineMetadata) error
	RemovePipeline(name string) error
	FindPipelineByName(name string) (PipelineMetadata, error)
	FindPipelines() ([]PipelineMetadata, error)

	SaveJob(job Job) error
	FindJob(groupId string) ([]Job, error)
	FindJobs(pipelineName string, query QueryParamJob) ([][]Job, error)
}

type PipelineMetadata struct {
	Name                string
	PathLocalRepository git.Path
	ManifestLocal       ServiceManifestLocal
	Service             MetadataService
}

type MetadataService struct {
	PathSrcDir                string
	PathBinDir                string
	PathEtcDir                string
	PathOptDir                string
	PathSystemdUnitFileDir    string
	PathSystemdUnitEnvFileDir string
	PathBackupDir             string
}
