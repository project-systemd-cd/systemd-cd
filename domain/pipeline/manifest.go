package pipeline

type ServiceManifestLocal struct {
	Name string `toml:"name"`

	GitRemoteUrl    string  `toml:"git_remote_url"`
	GitTargetBranch string  `toml:"git_target_branch"`
	GitTagRegex     *string `toml:"git_tag_regex"`
	GitManifestFile *string `toml:"git_manifest_file,omitempty"`

	TestCommands   *[]string       `toml:"test_commands"`
	BuildCommands  *[]string       `toml:"build_commands"`
	Binaries       *[]string       `toml:"binaries"`
	SystemdOptions []SystemdOption `toml:"systemd_services"`
}

type ServiceManifestRemote struct {
	Name           string          `toml:"name"`
	TestCommands   *[]string       `toml:"test_commands"`
	BuildCommands  *[]string       `toml:"build_commands"`
	Binaries       *[]string       `toml:"binaries"`
	SystemdOptions []SystemdOption `toml:"systemd_services"`
}

type ServiceManifestMerged struct {
	Name                  string                       `toml:"name"`
	GitTargetBranch       string                       `toml:"git_target_branch"`
	GitTagRegex           *string                      `toml:"git_tag_regex"`
	TestCommands          *[]string                    `toml:"test_commands"`
	BuildCommands         *[]string                    `toml:"build_commands"`
	Binaries              *[]string                    `toml:"binaries"`
	SystemdServiceOptions []SystemdServiceOptionMerged `toml:"systemd_services"`
}

type SystemdOption struct {
	Name         string       `toml:"name"`
	Description  *string      `toml:"description,omitempty"`
	ExecStartPre *string      `toml:"exec_start_pre,omitempty"`
	ExecStart    string       `toml:"exec_start"`
	Args         string       `toml:"args"`
	EnvVars      []EnvVar     `toml:"env"`
	Etc          []PathOption `toml:"etc,omitempty"`
	Opt          []string     `toml:"opt_files,omitempty"`
	Port         *uint16      `toml:"port,omitempty"`
}

type SystemdServiceOptionMerged struct {
	Name         string       `toml:"name"`
	Description  string       `toml:"description,omitempty"`
	ExecStartPre *string      `toml:"exec_start_pre,omitempty"`
	ExecStart    string       `toml:"exec_start"`
	Args         string       `toml:"args"`
	EnvVars      []EnvVar     `toml:"env"`
	Etc          []PathOption `toml:"etc,omitempty"`
	Opt          []string     `toml:"opt_files,omitempty"`
	Port         *uint16      `toml:"port,omitempty"`
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
