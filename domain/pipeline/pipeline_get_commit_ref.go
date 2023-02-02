package pipeline

func (p pipeline) GetCommitRef() string {
	return p.RepositoryLocal.RefCommitId
}
