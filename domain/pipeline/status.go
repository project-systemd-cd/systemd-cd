package pipeline

type (
	// Status of latest pipeline.
	// (Synced, Syncing, OutOfSync, Error)
	Status string
)

const (
	StatusSynced    Status = "synced"
	StatusSyncing   Status = "syncing"
	StatusOutOfSync Status = "out of sync"
	StatusError     Status = "error"
)
