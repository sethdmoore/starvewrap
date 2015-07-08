package signals

const (
	// shutdown must be zero as a closed channel emits zero-type
	SHUTDOWN = iota
	RESTART
	SAVE
	SIGINT
)
