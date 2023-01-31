package git

func (r *RepositoryLocal) Fetch() error {
	return r.git.command.Fetch(r.Path)
}
