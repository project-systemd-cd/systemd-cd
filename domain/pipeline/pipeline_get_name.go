package pipeline

func (p *pipeline) GetName() string {
	return p.ManifestMerged.Name
}
