package pipeline

import (
	"fmt"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"
)

func (m *ServiceManifestMerged) Validate() (err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Validate manifest")
	logger.Logger().Debugf("* manifestMerged.Name = %v", m.Name)
	logger.Logger().Tracef("* manifestMerged = %+v", *m)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Debug("END   - Validate manifest")
		} else {
			logger.Logger().Error("FAILED - Validate manifest")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	if m.Name == "" {
		err = &errors.ErrValidationMsg{Msg: "failed to validate manifest: 'name' is require"}
		return err
	}
	if m.TestCommands != nil {
		for i, cmd := range *m.TestCommands {
			if cmd == "" {
				err = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'test_commands[%d]' cannot be empty text", i)}
				return err
			}
		}
	}
	if m.BuildCommands != nil {
		for i, cmd := range *m.BuildCommands {
			if cmd == "" {
				err = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'build_commands[%d]' cannot be empty text", i)}
				return err
			}
		}
	}

	if m.Binaries != nil {
		for i, binary := range *m.Binaries {
			if binary == "" {
				err = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'binaries[%d]' cannot be empty text", i)}
				return err
			}
		}
	}
	for i, s := range m.SystemdOptions {
		// TODO: validate name duplication
		if s.Name == "" {
			err = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'systemd[%d].name' cannot be empty text", i)}
			return err
		}
		if s.ExecStart == "" {
			err = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'systemd[%d].exec_start' cannot be empty text", i)}
			return err
		}
		for j, etc := range s.Etc {
			if etc.Target == "" {
				err = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'systemd[%d].etc[%d].target' cannot be empty text", i, j)}
				return err
			}
			if etc.Option == "" {
				err = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'systemd[%d].etc[%d].option' cannot be empty text", i, j)}
				return err
			}
		}
		for j, opt := range s.Opt {
			if opt == "" {
				err = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'systemd[%d].opt_files[%d]' cannot be empty text", i, j)}
				return err
			}
		}
		// TODO: validate port duplication
		if s.Port != nil && (*s.Port < 1 || *s.Port > 65535) {
			err = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'systemd[%d].port' must be between 0 and 65535", i)}
			return err
		}
	}

	return nil
}
