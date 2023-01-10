package git

func (r *RepositoryLocal) DiffExists(executeFetch bool) (exists bool, err error) {
	if executeFetch {
		err = r.fetch()
		if err != nil {
			return
		}
	}
	return r.git.command.DiffExists(r.Path, r.TargetBranch)
}
