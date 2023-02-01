package git

func (r *RepositoryLocal) Checkout(hash string) (err error) {
	err = r.git.command.CheckoutHash(r.Path, hash)
	if err != nil {
		return
	}

	r.RefCommitId = hash

	return
}
