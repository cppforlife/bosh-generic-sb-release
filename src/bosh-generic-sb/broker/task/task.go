package task

type Func func() (value interface{}, err error)
type CancelFunc func(task Task) error

type State string

const (
	StateRunning State = "running"
	StateDone    State = "done"
	StateFailed  State = "failed"
)

type Task struct {
	ID    string
	State State

	Value interface{}
	Error error

	Func       Func
	CancelFunc CancelFunc
}

func (t Task) Cancel() error {
	if t.CancelFunc != nil {
		return t.CancelFunc(t)
	}
	return nil
}
