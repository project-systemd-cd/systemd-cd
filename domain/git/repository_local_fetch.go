package git

func (r *RepositoryLocal) fetch() error {
	return r.git.command.Fetch(r.Path)
}
