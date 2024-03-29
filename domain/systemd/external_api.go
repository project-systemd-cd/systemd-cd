package systemd

type Systemctl interface {
	DaemonReload() error
	Enable(service string, startNow bool) error
	Disable(service string, stopNow bool) error
	Start(service string) error
	Stop(service string) error
	Restart(service string) error
	Status(service string) (Status, error)
}
