package pipeline

func (p *pipeline) GetJobs(q QueryParamJob) ([][]Job, error) {
	return p.service.repo.FindJobs(p.ManifestMerged.Name, q)
}
