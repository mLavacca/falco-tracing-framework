package main

import (
	"fmt"
)

type Reporter struct {
}

func NewReporter(reportConf ReportConfiguration) *Reporter {
	r := new(Reporter)

	return r
}

func (r *Reporter) startReport() {
	fmt.Println("report")
}
