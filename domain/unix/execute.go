package unix

import (
	"bytes"
	"os/exec"
	"strings"
)

type ExecuteOption struct {
	WorkingDirectory *string
}

func Execute(o ExecuteOption, name string, arg ...string) (exitCode int, stdout bytes.Buffer, stderr bytes.Buffer, err error) {
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
		return
	}

	return
}
