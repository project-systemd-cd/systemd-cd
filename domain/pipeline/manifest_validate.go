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
