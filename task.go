package task

// A task is a set of commands that are run in a certain order
// and that may have errors and side effects
type Command func() (Command, error)

type Task interface {
	Run() error
}

func New(cmd Command) Task {
	return &task{nextCommand: cmd, err: nil}
}

type task struct {
	nextCommand Command
	err         error
}

func (t *task) Run() error {
	for t.nextCommand != nil {

		t.nextCommand, t.err = t.nextCommand()
	}

	return t.err
}
