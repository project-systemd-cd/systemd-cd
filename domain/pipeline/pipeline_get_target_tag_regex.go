package pipeline

func (p *pipeline) GetGitTargetTagRegex() *string {
	return p.ManifestLocal.GitTagRegex
}
