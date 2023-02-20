package pipeline

func (p pipeline) GetCommitMessage() (string, error) {
	return p.RepositoryLocal.GetCommitMessage(p.GetCommitRef())
}
