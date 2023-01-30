package git

func (r *RepositoryLocal) FindHashByTagRegex(regex string) (hash string, err error) {
	return r.git.command.FindHashByTagRegex(r.Path, regex)
}
