package systemctl

import (
	"bytes"
	"os/exec"
)

func executeCommand(name string, arg ...string) (exitCode int, stdout bytes.Buffer, stderr bytes.Buffer, err error) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	exitCode = cmd.ProcessState.ExitCode()
	return
}
