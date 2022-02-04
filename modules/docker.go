package modules

import (
	"os"
)

func init() {
	addModule("docketr", func() Module {
		return &docker{
			Base: Get(),
		}
	})
}

type docker struct {
	*Base

	// Output
	MachineName string
	Host        string
}

func (d *docker) Init() error {
	d.MachineName, _ = os.LookupEnv("DOCKER_MACHINE_NAME")
	d.Host, _ = os.LookupEnv("DOCKER_HOST")

	return nil
}
