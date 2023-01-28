package pipeline

type ServiceManifestLocal struct {
	GitRemoteUrl    string  `toml:"git_remote_url"`
	GitTargetBranch string  `toml:"git_target_branch"`
	GitManifestFile *string `toml:"git_manifest_file,omitempty"`

	Name           string        `toml:"name"`
	Description    *string       `toml:"description,omitempty"`
	Port           *uint16       `toml:"port,omitempty"`
	TestCommands   *[]string     `toml:"test_commands"`
	BuildCommands  *[]string     `toml:"build_commands"`
	Opt            *[]string     `toml:"opt_files"`
	Etc            *[]PathOption `toml:"etc"`
	EnvVars        []EnvVar      `toml:"env"`
	Binaries       *[]string     `toml:"binaries"`
	ExecuteCommand *string       `toml:"execute_command,omitempty"`
	Args           *string       `toml:"args,omitempty"`
}

type ServiceManifestRemote struct {
	Name          string       `toml:"name"`
	Description   string       `toml:"description"`
	Port          *uint16      `toml:"port,omitempty"`
	TestCommands  *[]string    `toml:"test_commands"`
	BuildCommands *[]string    `toml:"build_commands"`
	Opt           []string     `toml:"opt_files,omitempty"`
	Etc           []PathOption `toml:"etc,omitempty"`
	Binaries      *[]string    `toml:"binaries"`
	// TODO: omitempty
	ExecuteCommand string `toml:"execute_command"`
	Args           string `toml:"args"`
}

type ServiceManifestMerged struct {
	Name          string       `toml:"name"`
	Description   string       `toml:"description"`
	Port          *uint16      `toml:"port,omitempty"`
	TestCommands  *[]string    `toml:"test_commands"`
	BuildCommands *[]string    `toml:"build_commands"`
	Opt           []string     `toml:"opt_files,omitempty"`
	Etc           []PathOption `toml:"etc,omitempty"`
	EnvVars       []EnvVar     `toml:"env"`
	Binaries      *[]string    `toml:"binaries"`
	// TODO: omitempty
	ExecuteCommand string `toml:"execute_command"`
	Args           string `toml:"args"`
}

type PathOption struct {
	Target  string  `toml:"target"`
	Content *string `toml:"content,omitempty"`
	Option  string  `toml:"option"`
}

type EnvVar struct {
	Name  string `toml:"name"`
	Value string `toml:"value"`
}
