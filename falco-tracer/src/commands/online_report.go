package commands

import (
	"configuration"
	"log"
	"os"
	"os/exec"
	"stats_getter"
	"sync"
	"time"
)

type onlineReporter struct {
	reporter     reporterData
	pullInterval time.Duration
}

func newOnlineReporter(conf configuration.OnlineReportConfiguration) *onlineReporter {
	r := new(onlineReporter)

	r.reporter.falcoBins = conf.ProgConfig.ProgBins
	r.reporter.falcoargs = conf.ProgConfig.ProgArgs
	r.reporter.outputDirectory = conf.OutputDirectory
	r.reporter.outputFormats = conf.OutputFormats
	r.reporter.mode = "online"

	r.pullInterval = conf.PullInterval

	return r
}

func (or *onlineReporter) report() {
	or.startReport()
}

func (r *onlineReporter) startReport() {

	for _, bin := range r.reporter.falcoBins {
		args := r.reporter.falcoargs

		cmd := exec.Command(bin, args...)
		err := cmd.Start()
		if err != nil {
			log.Fatalln("cmd.Start() failed with ", err)
		}

		r.reporter.falcoTracer = stats_getter.NewFalcoTracer(r.reporter.mode)

		var wg sync.WaitGroup
		sigs := make(chan os.Signal)

		r.reporter.falcoTracer.LoadOnlineRulesFromFalco()

		r.reporter.falcoTracer.FlushFalcoData()

		wg.Add(1)
		go r.reporter.falcoTracer.LoadOnlineStatsFromFalco(r.pullInterval, &wg)

		<-sigs

		r.reporter.falcoTracer.ExitFlag = true

		wg.Wait()

		//jsonStats, err := r.reporter.falcoTracer.MarshalJSON()
		if err != nil {
			log.Fatal("error in object marshaling")
		}

		//writeMetricsOnFile(jsonStats, r.reporter.outputFile)
	}

}
