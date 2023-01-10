package git

func (r *RepositoryLocal) Pull(force bool) (refCommitId string, err error) {
	refCommitId, err = r.git.command.Pull(r.Path, force)
	if err != nil {
		return
	}
	r.RefCommitId = refCommitId
	return
}
