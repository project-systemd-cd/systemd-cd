package git

type (
	Path string

	RepositoryLocal struct {
		git          *Git
		RemoteUrl    string `toml:"git_remote_url"`
		TargetBranch string `toml:"target_branch"`
		RefCommitId  string `toml:"ref_commit_id"`
		Path         Path   `toml:"path"`
	}
)
