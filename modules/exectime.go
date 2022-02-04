package modules

import "fmt"

func init() {
	addModule("execTime", func() Module {
		return &execTime{
			Base: Get(),
		}
	})
}

type execTime struct {
	*Base

	// Output
	Time  int64
	Value string
}

func (t *execTime) Init() error {
	t.Time = int64(t.Base.executionTime)
	t.Value = fmt.Sprintf("%dms", t.Time%1000)
	return nil
}
