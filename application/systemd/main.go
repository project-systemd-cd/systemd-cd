package systemd

import "systemd-cd/domain/model/systemd"

func New(s systemd.Systemctl, unitFileDir string) (systemd.ISystemd, error) {
	return systemd.New(s, unitFileDir)
}

type (
	UnitFileService  = systemd.UnitFileService
	UnitDirective    = systemd.UnitDirective
	UnitType         = systemd.UnitType
	ServiceDirective = systemd.ServiceDirective
	InstallDirective = systemd.InstallDirective
)

var (
	UnitTypeSimple  = systemd.UnitTypeSimple
	UnitTypeForking = systemd.UnitTypeForking
	UnitTypeOneShot = systemd.UnitTypeOneShot
	UnitTypeDbus    = systemd.UnitTypeDbus
	UnitTypeNotify  = systemd.UnitTypeNotify
	UnitTypeIdle    = systemd.UnitTypeIdle
)
