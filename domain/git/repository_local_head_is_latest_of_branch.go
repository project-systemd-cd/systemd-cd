package git

func (r *RepositoryLocal) HeadIsLatesetOfBranch(branch string) (bool, error) {
	exists, err := r.git.command.DiffExists(r.Path, branch)
	if err != nil {
		return false, err
	}
	isLatest := !exists

	return isLatest, nil
}
