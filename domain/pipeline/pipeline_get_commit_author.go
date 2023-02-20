package pipeline

func (p pipeline) GetCommitAuthor() (string, error) {
	return p.RepositoryLocal.GetCommitAuthor(p.GetCommitRef())
}
