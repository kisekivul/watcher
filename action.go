package watcher

// Action action of path
type Action int

const (
	// NONE none action
	NONE Action = iota
	// CREATE create action
	CREATE
	// UPDATE update action
	UPDATE
	// REMOVE remove action
	REMOVE
)

var (
	actions = [...]string{
		"none",
		"create",
		"update",
		"remove",
	}
)

// String return actions string
func (a Action) String() string {
	return actions[int(a)%len(actions)]
}
