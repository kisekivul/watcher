package watcher

// Event event of path
type Event struct {
	Date   int64
	Name   string
	Kind   Kind
	Action Action
}
