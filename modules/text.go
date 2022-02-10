package modules

func init() {
	addModule("text", func() Module {
		return &text{
			Base: Get(),
		}
	})
}

type text struct {
	*Base
}

func (t *text) Init() error {
	return nil
}
