package systemctl

import (
	"systemd-cd/domain/systemd"
)

func New() systemd.Systemctl {
	return systemctl{}
}

type systemctl struct{}
