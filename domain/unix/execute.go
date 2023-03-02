package unix

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"systemd-cd/domain/logger"
)

type ExecuteOption struct {
	WorkingDirectory *string
	WantExitCodes    []int
}

func Execute(o ExecuteOption, name string, arg ...string) (exitCode int, stdout bytes.Buffer, stderr bytes.Buffer, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Execute command")
	logger.Logger().Debugf("< command = %v", strings.Join(append([]string{name}, arg...), " "))
	logger.Logger().Debugf("< option = %+v", o)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil {
			logger.Logger().Tracef("> stdout = %+v", stdout.String())
			logger.Logger().Tracef("> exit code = %+v", exitCode)
			logger.Logger().Debug("END   - Execute command")
		} else {
			logger.Logger().Error("FAILED - Execute command")
			logger.Logger().Errorf("command = %v", strings.Join(append([]string{name}, arg...), " "))
			logger.Logger().Error(err)
			logger.Logger().Error("exit status", exitCode)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	containWildcard := false
	for _, a := range arg {
		if strings.Contains(a, "*") {
			containWildcard = true
			break
		}
	}

	if o.WorkingDirectory != nil {
		command := strings.Join(append([]string{"cd", *o.WorkingDirectory, ";", name}, arg...), " ")
		name = "/usr/bin/bash"
		arg = []string{"-c", command}
	} else if containWildcard {
		command := strings.Join(append([]string{name}, arg...), " ")
		name = "/usr/bin/bash"
		arg = []string{"-c", command}
	}

	cmd := exec.Command(name, arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	exitCode = cmd.ProcessState.ExitCode()
	if err != nil {
		for _, code := range o.WantExitCodes {
			if exitCode == code {
				err = nil
				return
			}
		}
		if strings.HasPrefix(err.Error(), "exit status") {
			err = errors.New(stderr.String())
		}
		return
	}
	return
}
