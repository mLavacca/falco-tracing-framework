package commands

import (
	"configuration"
	"fmt"
)

type Reporter struct {
}

func NewReporter(conf configuration.TracerConfigurations) *Reporter {
	r := new(Reporter)

	return r
}

func (r *Reporter) StartReport() {
	fmt.Println("report")
}
