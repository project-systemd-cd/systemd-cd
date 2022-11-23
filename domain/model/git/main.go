package git

func New(git GitCommand) *Git {
	return &Git{command: git}
}

type Git struct {
	command GitCommand
}
