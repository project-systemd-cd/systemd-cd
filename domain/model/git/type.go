package git

type (
	Path string

	RepositoryRemote struct {
		RemoteUrl string  `toml:"git_remote_url"`
		User      *string `toml:"git_user,omitempty"`
		Token     *string `toml:"git_access_token,omitempty"`
	}

	RepositoryLocal struct {
		git *Git
		RepositoryRemote
		TargetBranch string `toml:"target_branch"`
		RefCommitId  string `toml:"ref_commit_id"`
		Path         Path   `toml:"path"`
	}
)
