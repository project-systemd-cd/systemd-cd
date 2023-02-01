package pipeline

import (
	"systemd-cd/domain/unix"
)

func (p pipeline) build() error {
	if p.ManifestMerged.BuildCommands != nil {
		for _, cmd := range *p.ManifestMerged.BuildCommands {
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
