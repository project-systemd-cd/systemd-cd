package toml

import (
	"systemd-cd/domain/unix"
)

// RemovePipeline implements pipeline.IRepository
func (r *rPipeline) RemovePipeline(name string) error {
	_, err := r.FindPipelineByName(name)
	if err != nil {
		return err
	}

	// Remove files
	err = unix.Rm(unix.ExecuteOption{WantExitCodes: []int{1}}, unix.RmOption{}, r.basePath+name+".toml")
	if err != nil {
		return err
	}

	return nil
}
