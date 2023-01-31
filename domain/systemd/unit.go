package systemd

type Unit interface {
	Enable(startNow bool) error
	Disable(stopNow bool) error
	Start() error
	Stop() error
	Restart() error
	GetStatus() (Status, error)
}

type (
	Status string
)

const (
	// Systemd service status
	StatusStopped Status = "stopped"
	StatusRunning Status = "running"
	StatusFailed  Status = "failed"
)

var (
	// check implements
	_ Unit = UnitService{}
)

type (
	// +Unit
	UnitService struct {
		systemctl             Systemctl
		Name                  string
		unitFile              UnitFileService
		Path                  string
		EnvironmentFileValues map[string]string
	}
)

// +Unit
func (u UnitService) Enable(startNow bool) error {
	return u.systemctl.Enable(u.Name, startNow)
}

// +Unit
func (u UnitService) Disable(stopNow bool) error {
	return u.systemctl.Disable(u.Name, stopNow)
}

// +Unit
func (u UnitService) Start() error {
	return u.systemctl.Start(u.Name)
}

// +Unit
func (u UnitService) Stop() error {
	return u.systemctl.Stop(u.Name)
}

// +Unit
func (u UnitService) Restart() error {
	return u.systemctl.Restart(u.Name)
}

// +Unit
func (u UnitService) GetStatus() (Status, error) {
	return u.systemctl.Status(u.Name)
}
