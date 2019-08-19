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

	duration int
	tester   *falco_test.Tester
}

func newRecorder(conf configuration.TracerConfigurations) *Recorder {
	r := new(Recorder)
	var d int

	r.sysdigBin = conf.Record.ProgConfig.ProgBin
	r.sysdigArgs = conf.Record.ProgConfig.ProgArgs
	r.tester, d = falco_test.NewTester(conf)
	if r.tester == nil {
		r.duration = d
	}

	return r
}

func (r *Recorder) startRecord() {

	cmd := exec.Command(r.sysdigBin, r.sysdigArgs...)

	err := cmd.Start()
	if err != nil {
		log.Fatalln("cmd.Start() failed with ", err)
	}

	if r.tester != nil {
		time.Sleep(2 * time.Second)
		r.tester.RunAllTests()
		time.Sleep(2 * time.Second)
	} else {
		time.Sleep(time.Duration(r.duration) * time.Second)
	}

	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		log.Fatalln("failed to kill process: ", err)
	}

	cmd.Wait()
}

func (r *Recorder) rollback() {
	if r.tester != nil {
		r.tester.RunAllRollbacks()
	}
}

func (r *Recorder) Record() {
	r.startRecord()
	r.rollback()
}
