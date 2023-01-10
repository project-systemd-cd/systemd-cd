package systemd

import "os"

// DeleteService implements iSystemdService
func (s Systemd) DeleteService(u UnitService) error {
	err := u.Disable(true)
	if err != nil {
		return err
	}

	// Delete `.service` file
	err = os.Remove(u.Path)

	return err
}
