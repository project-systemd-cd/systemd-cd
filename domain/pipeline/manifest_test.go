package pipeline_test

import (
	"bytes"
	"reflect"
	"systemd-cd/domain/pipeline"
	"systemd-cd/domain/toml"
	"testing"
)

func TestUnmarshalManifest(t *testing.T) {
	const manifest = `name = "systemd-cd"
git_remote_url = "https://github.com/tingtt/systemd-cd.git"
git_target_branch = "main"
git_tag_regex = "v*"
build_commands = ["/usr/bin/go build"]
binaries = ["systemd-cd"]

[[systemd]]
name = "systemd-cd"
description = "systemd-cd"
exec_start_pre = "/usr/bin/go version"
exec_start = "systemd-cd"
args = "--log.level debug"
opt_files = ["README.md"]
port = 443`

	type args struct {
		ManifestContent string
	}
	tests := []struct {
		name string
		args args
		want pipeline.ServiceManifestLocal
	}{{
		name: "equals",
		args: args{
			ManifestContent: manifest,
		},
		want: pipeline.ServiceManifestLocal{
			Name:            "systemd-cd",
			GitRemoteUrl:    "https://github.com/tingtt/systemd-cd.git",
			GitTargetBranch: "main",
			GitTagRegex:     func() *string { s := "v*"; return &s }(),
			GitManifestFile: nil,
			TestCommands:    nil,
			BuildCommands:   &[]string{"/usr/bin/go build"},
			Binaries:        &[]string{"systemd-cd"},
			SystemdOptions: []pipeline.SystemdOption{{
				Name:         "systemd-cd",
				Description:  func() *string { s := "systemd-cd"; return &s }(),
				ExecStartPre: func() *string { s := "/usr/bin/go version"; return &s }(),
				ExecStart:    "systemd-cd",
				Args:         "--log.level debug",
				EnvVars:      nil,
				Etc:          nil,
				Opt:          []string{"README.md"},
				Port:         func() *uint16 { s := uint16(443); return &s }(),
			}},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &bytes.Buffer{}
			b.WriteString(tt.args.ManifestContent)
			m := new(pipeline.ServiceManifestLocal)
			err := toml.Decode(b, m)
			if err != nil {
				t.Errorf("toml.Decode() error = %v", err)
				return
			}
			if !reflect.DeepEqual(*m, tt.want) {
				t.Errorf("toml.Decode() = %+v, want %+v", *m, tt.want)
			}
		})
	}
}
