package main

import (
	"log"
	"os/exec"
)

type Recorder struct {
	sysdigBin  string
	sysdigArgs []string

	tester Tester
}

func NewRecorder(recordConf RecordConfiguration) *Recorder {
	r := new(Recorder)

	r.sysdigBin = recordConf.ProgConfig.ProgBin
	r.sysdigArgs = recordConf.ProgConfig.ProgArgs

	return r
}

func (r *Recorder) startRecord() {

	cmd := exec.Command(r.sysdigBin, r.sysdigArgs...)

	err := cmd.Run()
	if err != nil {
		log.Fatalln("cmd.Run() failed with ", err)
	}
}
