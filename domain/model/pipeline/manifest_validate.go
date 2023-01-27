package pipeline

import (
	"errors"
	"systemd-cd/domain/model/logger"
)

func (m *ServiceManifestMerged) Validate() error {
	logger.Logger().Trace(logger.Var2Text(
		"Called",
		[]logger.Var{
			{Value: m},
		},
	))

	if m.Name == "" {
		err := errors.New("failed to validate manifest: 'name' is require")
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}
	if m.Port != nil && (*m.Port < 1 || *m.Port > 65535) {
		err := errors.New("failed to validate manifest: 'port' must be between 0 and 65535")
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}
	if m.TestCommand != nil && *m.TestCommand == "" {
		err := errors.New("failed to validate manifest: 'test_command' cannot be empty")
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}
	if m.BuildCommand != nil && *m.BuildCommand == "" {
		err := errors.New("failed to validate manifest: 'build_command' cannot be empty")
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}
	if m.Binary != nil && *m.Binary == "" {
		err := errors.New("failed to validate manifest: 'binary' cannot be empty")
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}
	if m.ExecuteCommand == "" {
		err := errors.New("failed to validate manifest: 'execute_command' is require")
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return err
	}

	logger.Logger().Trace("Finished")
	return nil
}
