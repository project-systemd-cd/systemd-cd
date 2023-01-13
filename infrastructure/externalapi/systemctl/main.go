package systemctl

import (
	"systemd-cd/domain/model/systemd"
)

func New() systemd.Systemctl {
	return systemctl{}
}

type systemctl struct{}
