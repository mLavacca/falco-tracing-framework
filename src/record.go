package main

type Recorder struct {
	sysdigBin  string
	sysdigArgs []string

	tester *Tester
}

func NewRecorder(conf TracerConfigurations) *Recorder {
	r := new(Recorder)

	r.sysdigBin = conf.Record.ProgConfig.ProgBin
	r.sysdigArgs = conf.Record.ProgConfig.ProgArgs
	r.tester = NewTester(conf)

	return r
}

func (r *Recorder) startRecord() {

	/*cmd := exec.Command(r.sysdigBin, r.sysdigArgs...)

	err := cmd.Run()
	if err != nil {
		log.Fatalln("cmd.Run() failed with ", err)
	}
	*/
	r.tester.runAllTests()
}
