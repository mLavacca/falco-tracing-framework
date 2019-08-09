package commands

import (
	"configuration"
	"falco_test"
	"log"
	"os/exec"
	"syscall"
	"time"
)

type Recorder struct {
	sysdigBin  string
	sysdigArgs []string

	tester *falco_test.Tester
}

func NewRecorder(conf configuration.TracerConfigurations) *Recorder {
	r := new(Recorder)

	r.sysdigBin = conf.Record.ProgConfig.ProgBin
	r.sysdigArgs = conf.Record.ProgConfig.ProgArgs
	r.tester = falco_test.NewTester(conf)

	return r
}

func (r *Recorder) StartRecord() {

	cmd := exec.Command(r.sysdigBin, r.sysdigArgs...)

	err := cmd.Start()
	if err != nil {
		log.Fatalln("cmd.Run() failed with ", err)
	}

	time.Sleep(3 * time.Second)

	r.tester.RunAllTests()

	time.Sleep(3 * time.Second)

	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		log.Fatalln("failed to kill process: ", err)
	}
}

func (r *Recorder) Rollback() {
	r.tester.RunAllRollbacks()
}
