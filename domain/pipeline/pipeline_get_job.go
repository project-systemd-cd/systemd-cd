package pipeline

func (p *pipeline) GetJob(groupId string) ([]Job, error) {
	return p.service.repo.FindJob(groupId)
}
