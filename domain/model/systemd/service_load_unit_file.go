package systemd

import (
	"bytes"
	"strings"
)

// loadUnitFileSerivce implements iSystemdService
func (s Systemd) loadUnitFileSerivce(path string) (u UnitFileService, isGeneratedBySystemdCd bool, err error) {
	// Read file
	b := &bytes.Buffer{}
	err = readFile(path, b)
	if err != nil {
		return
	}

	// Check generator
	if strings.Contains(b.String(), "#! Generated by systemd-cd\n") {
		isGeneratedBySystemdCd = true
	}

	// Unmarshal
	u, err = UnmarshalUnitFile(b)

	return
}