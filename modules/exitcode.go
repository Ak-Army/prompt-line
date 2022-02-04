package modules

func init() {
	addModule("exitCode", func() Module {
		return &exitCode{
			Base: Get(),
		}
	})
}

var exitCodes = map[int]string{
	1:   "ERROR",
	2:   "USAGE",
	127: "NOTFOUND",
}

type exitCode struct {
	*Base

	// Output
	Code int
	Name string
}

func (t *exitCode) Init() error {
	t.Code = t.Base.exitCode
	if v, ok := exitCodes[t.Code]; ok {
		t.Name = v
	}
	return nil
}
