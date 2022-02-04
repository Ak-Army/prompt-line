package modules

import "time"

func init() {
	addModule("dateTime", func() Module {
		return &dateTime{
			Base: Get(),
		}
	})
}

type dateTime struct {
	*Base

	// Output
	Time time.Time
}

func (t *dateTime) Init() error {
	t.Time = time.Now()

	return nil
}
