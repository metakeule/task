package task

// Task is a function with sideeffects that returns
// the next task to be done, after it has been done
// its job. If there is no task to proceed, it returns nil instead.
// If there is an error resulting from running its job,
// the error is returned and the returned task must be nil as well.
type Task func() (Task, error)

// Run calls the given tasks and any of the following
// tasks until an error happened of the task is nil.
// The first error is returned.
func Run(t Task) (err error) {
	for t != nil && err == nil {
		t, err = t()
	}

	return
}
