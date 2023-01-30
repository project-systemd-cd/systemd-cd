package pipeline

// TODO: .Opt -> .SystemdOptions[i].Opt
type ServiceManifestLocal struct {
	GitRemoteUrl    string  `toml:"git_remote_url"`
	GitTargetBranch string  `toml:"git_target_branch"`
	GitManifestFile *string `toml:"git_manifest_file,omitempty"`

	Name           string          `toml:"name"`
	TestCommands   *[]string       `toml:"test_commands"`
	BuildCommands  *[]string       `toml:"build_commands"`
	Opt            *[]string       `toml:"opt_files"`
	Binaries       *[]string       `toml:"binaries"`
	SystemdOptions []SystemdOption `toml:"systemd"`
}

// TODO: .Opt -> .SystemdOptions[i].Opt
type ServiceManifestRemote struct {
	Name           string          `toml:"name"`
	TestCommands   *[]string       `toml:"test_commands"`
	BuildCommands  *[]string       `toml:"build_commands"`
	Opt            []string        `toml:"opt_files,omitempty"`
	Binaries       *[]string       `toml:"binaries"`
	SystemdOptions []SystemdOption `toml:"systemd"`
}

// TODO: .Opt -> .SystemdOptions[i].Opt
type ServiceManifestMerged struct {
	Name           string                `toml:"name"`
	TestCommands   *[]string             `toml:"test_commands"`
	BuildCommands  *[]string             `toml:"build_commands"`
	Opt            []string              `toml:"opt_files,omitempty"`
	Binaries       *[]string             `toml:"binaries"`
	SystemdOptions []SystemdOptionMerged `toml:"systemd"`
}

type SystemdOption struct {
	Name           string       `toml:"name"`
	Description    *string      `toml:"description,omitempty"`
	ExecuteCommand string       `toml:"execute_command"`
	Args           string       `toml:"args"`
	EnvVars        []EnvVar     `toml:"env"`
	Etc            []PathOption `toml:"etc,omitempty"`
	Port           *uint16      `toml:"port,omitempty"`
}

type SystemdOptionMerged struct {
	Name           string       `toml:"name"`
	Description    string       `toml:"description,omitempty"`
	ExecuteCommand string       `toml:"execute_command"`
	Args           string       `toml:"args"`
	EnvVars        []EnvVar     `toml:"env"`
	Etc            []PathOption `toml:"etc,omitempty"`
	Port           *uint16      `toml:"port,omitempty"`
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
