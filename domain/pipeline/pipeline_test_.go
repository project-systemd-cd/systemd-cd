package pipeline

import (
	"systemd-cd/domain/unix"
)

func (p pipeline) test() error {
	if p.ManifestMerged.TestCommands != nil {
		for _, cmd := range *p.ManifestMerged.TestCommands {
			_, _, _, err := unix.Execute(
				unix.ExecuteOption{
					WorkingDirectory: (*string)(&p.RepositoryLocal.Path),
				},
				"/usr/bin/bash", "-c", "\""+cmd+"\"",
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
