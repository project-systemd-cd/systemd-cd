package pipeline

type (
	// Status of latest pipeline.
	// (Synced, Syncing, OutOfSync, Error)
	status string
)

const (
	StatusSynced    status = "synced"
	StatusSyncing   status = "syncing"
	StatusOutOfSync status = "out of sync"
	StatusFailed    status = "failed"
)
