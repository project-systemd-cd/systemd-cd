package systemd

import (
	"bytes"
	"strings"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/toml"
)

// loadEnvFile implements iSystemdService
func (s Systemd) loadEnvFile(path string) (e map[string]string, isGeneratedBySystemdCd bool, err error) {
	logger.Logger().Trace(logger.Var2Text("Called", []logger.Var{{Name: "path", Value: path}}))

	// Read file
	b := &bytes.Buffer{}
	err = readFile(path, b)
	if err != nil {
		logger.Logger().Error(logger.Var2Text("Error", []logger.Var{{Name: "err", Value: err}}))
		return
	}

	// Check generator
	if strings.Contains(b.String(), "#! Generated by systemd-cd\n") {
		isGeneratedBySystemdCd = true
	}

	// Decode
	// TODO: toml format で問題ないか検証
	logger.Logger().Warn("Unchecked code: no problem to systemd service environment file with TOML format.")
	err = toml.Decode(b, &e)
	if err != nil {
		logger.Logger().Error("Error:\n\tisGeneratedBySystemdCd: %v\n\terr: %v", isGeneratedBySystemdCd, err)
		return
	}

	logger.Logger().Trace(logger.Var2Text("Finished", []logger.Var{{Name: "isGeneratedBySystemdCd", Value: isGeneratedBySystemdCd}}))
	return
}