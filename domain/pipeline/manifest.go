package pipeline

type ServiceManifestLocal struct {
	GitRemoteUrl    string  `toml:"git_remote_url"`
	GitTargetBranch string  `toml:"git_target_branch"`
	GitManifestFile *string `toml:"git_manifest_file,omitempty"`

	Name           string        `toml:"name"`
	Description    *string       `toml:"description,omitempty"`
	Port           *uint16       `toml:"port,omitempty"`
	TestCommands   *[]string     `toml:"test_commands,omitempty"`
	BuildCommands  *[]string     `toml:"build_commands,omitempty"`
	Opt            *[]string     `toml:"src,omitempty"`
	Etc            *[]PathOption `toml:"etc,omitempty"`
	EnvVars        []EnvVar      `toml:"env_variables,omitempty"`
	Binaries       *[]string     `toml:"binaries,omitempty"`
	ExecuteCommand *string       `toml:"execute_command,omitempty"`
	Args           *string       `toml:"args,omitempty"`
}

type ServiceManifestRemote struct {
	Name          string       `toml:"name"`
	Description   string       `toml:"description"`
	Port          *uint16      `toml:"port,omitempty"`
	TestCommands  *[]string    `toml:"test_commands,omitempty"`
	BuildCommands *[]string    `toml:"build_commands,omitempty"`
	Opt           []string     `toml:"src,omitempty"`
	Etc           []PathOption `toml:"etc,omitempty"`
	Binaries      *[]string    `toml:"binaries,omitempty"`
	// TODO: omitempty
	ExecuteCommand string `toml:"execute_command"`
	Args           string `toml:"args"`
}

type ServiceManifestMerged struct {
	Name          string       `toml:"name"`
	Description   string       `toml:"description"`
	Port          *uint16      `toml:"port,omitempty"`
	TestCommands  *[]string    `toml:"test_commands,omitempty"`
	BuildCommands *[]string    `toml:"build_commands,omitempty"`
	Opt           []string     `toml:"src,omitempty"`
	Etc           []PathOption `toml:"etc,omitempty"`
	EnvVars       []EnvVar     `toml:"env_variables,omitempty"`
	Binaries      *[]string    `toml:"binaries,omitempty"`
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
