package globals

// Probably don't need this
const (
	OK = iota
	WARNING
	ERROR
)

// for checking triggers
const (
	_ = iota // index starts at 0 (unintialized)
	MATCH
	NO_MATCH
	HELP
)
