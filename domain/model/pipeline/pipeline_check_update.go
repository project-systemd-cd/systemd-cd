package pipeline

func (p *pipeline) CheckUpdate() (updateExists bool, err error) {
	// Check update
	outOfSync, err := p.RepositoryLocal.DiffExists(true)
	if outOfSync {
		p.Status = StatusOutOfSync
	} else {
		p.Status = StatusSynced
	}
	return outOfSync, err
}
