package unix

import (
	"bytes"
	"os/exec"
	"strings"
	"systemd-cd/domain/logger"
)

type ExecuteOption struct {
	WorkingDirectory *string
}

func Execute(o ExecuteOption, name string, arg ...string) (exitCode int, stdout bytes.Buffer, stderr bytes.Buffer, err error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Value: o}, {Name: "name", Value: name}, {Name: "arg", Value: arg}}))

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

	logger.Logger().Debugf("Debug:\n\tCommand: %v", strings.Join(append([]string{name}, arg...), " "))
	cmd := exec.Command(name, arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	exitCode = cmd.ProcessState.ExitCode()
	if err != nil {
		logger.Logger().Error("Error:\n\terr: %v\n\tstderr: %v", err, stderr.String())
		return
	}

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "stdout", Value: stdout.String()}}))
	return
}
