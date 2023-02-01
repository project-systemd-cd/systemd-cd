package pipeline

import (
	"fmt"
	"systemd-cd/domain/errors"
)

func (m *ServiceManifestMerged) Validate() error {
	if m.Name == "" {
		var err error = &errors.ErrValidationMsg{Msg: "failed to validate manifest: 'name' is require"}
		return err
	}
	if m.TestCommands != nil {
		for i, cmd := range *m.TestCommands {
			if cmd == "" {
				var err error = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'test_commands[%d]' cannot be empty text", i)}
				return err
			}
		}
	}
	if m.BuildCommands != nil {
		for i, cmd := range *m.BuildCommands {
			if cmd == "" {
				var err error = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'build_commands[%d]' cannot be empty text", i)}
				return err
			}
		}
	}

	if m.Binaries != nil {
		for i, binary := range *m.Binaries {
			if binary == "" {
				var err error = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'binaries[%d]' cannot be empty text", i)}
				return err
			}
		}
	}
	for i, s := range m.SystemdOptions {
		// TODO: validate name duplication
		if s.Name == "" {
			var err error = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'systemd[%d].name' cannot be empty text", i)}
			return err
		}
		if s.ExecuteCommand == "" {
			var err error = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'systemd[%d].execute_command' cannot be empty text", i)}
			return err
		}
		for j, etc := range s.Etc {
			if etc.Target == "" {
				var err error = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'systemd[%d].etc[%d].target' cannot be empty text", i, j)}
				return err
			}
			if etc.Option == "" {
				var err error = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'systemd[%d].etc[%d].option' cannot be empty text", i, j)}
				return err
			}
		}
		for j, opt := range s.Opt {
			if opt == "" {
				var err error = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'systemd[%d].opt_files[%d]' cannot be empty text", i, j)}
				return err
			}
		}
		// TODO: validate port duplication
		if s.Port != nil && (*s.Port < 1 || *s.Port > 65535) {
			var err error = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'systemd[%d].port' must be between 0 and 65535", i)}
			return err
		}
	}

	return nil
}
