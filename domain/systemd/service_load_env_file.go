package systemd

import (
	"bytes"
	errorss "errors"
	"os"
	"strings"
	"systemd-cd/domain/logger"
	"systemd-cd/domain/toml"
	"systemd-cd/domain/unix"
)

// loadEnvFile implements iSystemdService
func (s Systemd) loadEnvFile(path string) (e map[string]string, isGeneratedBySystemdCd bool, err error) {
	logger.Logger().Debug("-----------------------------------------------------------")
	logger.Logger().Debug("START - Load systemd env file")
	logger.Logger().Debugf("< path = %v", path)
	logger.Logger().Debug("-----------------------------------------------------------")
	defer func() {
		logger.Logger().Debug("-----------------------------------------------------------")
		if err == nil || errorss.Is(err, os.ErrNotExist) {
			logger.Logger().Debugf("> e = %+v", e)
			logger.Logger().Debugf("> isGeneratedBySystemdCd = %v", isGeneratedBySystemdCd)
			logger.Logger().Debug("END   - Load systemd env file")
		} else {
			logger.Logger().Error("FAILED - Load systemd env file")
			logger.Logger().Error(err)
		}
		logger.Logger().Debug("-----------------------------------------------------------")
	}()

	// Read file
	b := &bytes.Buffer{}
	err = unix.ReadFile(path, b)
	if err != nil {
		return
	}

	// Check generator
	if strings.Contains(b.String(), "#! Generated by systemd-cd\n") {
		isGeneratedBySystemdCd = true
	}

	// Decode
	err = toml.Decode(b, &e)
	if err != nil {
		return
	}

	return
}
