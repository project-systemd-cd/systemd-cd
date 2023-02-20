package pipeline

func (p *pipeline) GetGitTargetBranch() string {
	return p.ManifestLocal.GitTargetBranch
}
