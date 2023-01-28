package pipeline

import (
	"fmt"
	"systemd-cd/domain/errors"
	"systemd-cd/domain/logger"
)

func (m *ServiceManifestMerged) Validate() error {
	logger.Logger().Trace(logger.Var2Text(
		"Called",
		[]logger.Var{
			{Value: m},
		},
	))

	if m.Name == "" {
		var err error = &errors.ErrValidationMsg{Msg: "failed to validate manifest: 'name' is require"}
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}
	if m.Port != nil && (*m.Port < 1 || *m.Port > 65535) {
		var err error = &errors.ErrValidationMsg{Msg: "failed to validate manifest: 'port' must be between 0 and 65535"}
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}
	if m.TestCommands != nil {
		for i, cmd := range *m.TestCommands {
			if cmd == "" {
				var err error = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'test_commands[%d]' cannot be empty text", i)}
				logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
				return err
			}
		}
	}
	if m.BuildCommands != nil {
		for i, cmd := range *m.BuildCommands {
			if cmd == "" {
				var err error = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'build_commands[%d]' cannot be empty text", i)}
				logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
				return err
			}
		}
	}
	for i, opt := range m.Opt {
		if opt == "" {
			var err error = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'opt_files[%d]' cannot be empty text", i)}
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return err
		}
	}
	for i, etc := range m.Etc {
		if etc.Target == "" {
			var err error = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'etc[%d].target' cannot be empty text", i)}
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return err
		}
		if etc.Option == "" {
			var err error = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'etc[%d].option' cannot be empty text", i)}
			logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
			return err
		}
	}
	// TODO: validate env
	logger.Logger().Warn("Environment variables are unsafe because they are not checked for security")
	// for i, ev := range m.EnvVars {}
	if m.Binaries != nil {
		for i, binary := range *m.Binaries {
			if binary == "" {
				var err error = &errors.ErrValidationMsg{Msg: fmt.Sprintf("failed to validate manifest: 'binaries[%d]' cannot be empty text", i)}
				logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
				return err
			}
		}
	}
	// TODO: omitempty
	if m.ExecuteCommand == "" {
		var err error = &errors.ErrValidationMsg{Msg: "failed to validate manifest: 'execute_command' is require"}
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}
