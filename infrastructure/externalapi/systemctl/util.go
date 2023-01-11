package systemctl

import (
	"bytes"
	"os/exec"
	"systemd-cd/domain/model/logger"
)

func executeCommand(name string, arg ...string) (exitCode int, stdout bytes.Buffer, stderr bytes.Buffer, err error) {
	logger.Logger().Trace("Called:\n\targ.name: %v\n\targ.arg: %v", name, arg)

	cmd := exec.Command(name, arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	exitCode = cmd.ProcessState.ExitCode()

	if err != nil {
		logger.Logger().Errorf("Error:\n\terr: %v\n\tstdout: %v\n\tstderr: %v\n\texitCode: %v", err, stdout.String(), stderr.String(), exitCode)
		return
	}
	logger.Logger().Tracef("Finished:\n\tstdout: %v\n\tstderr: %v\n\texitCode: %v", stdout.String(), stderr.String(), exitCode)
	return
}
