package git

func (r *RepositoryLocal) CheckoutBranch(branch string) (err error) {
	err = r.git.command.CheckoutBranch(r.Path, branch)
	if err != nil {
		return
	}

	s, err := r.git.command.RefCommitId(r.Path)
	if err != nil {
		return
	}

	r.RefCommitId = s

	return
}
