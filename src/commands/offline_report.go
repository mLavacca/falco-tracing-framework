package commands

import (
	"configuration"
	"log"
	"os/exec"
	"stats_getter"
)

type offlineReporter struct {
	reporter reporterData

	iterations int
}

func newOfflineReporter(conf configuration.OfflineReportConfiguration) *offlineReporter {
	r := new(offlineReporter)

	r.reporter.falcoBin = conf.ProgConfig.ProgBin
	r.reporter.falcoargs = conf.ProgConfig.ProgArgs
	r.reporter.outputFile = conf.OutputFile
	r.reporter.mode = "offline"

	r.iterations = conf.Iterations

	return r
}

func (or *offlineReporter) report() {
	or.startReport()
}

func (r *offlineReporter) startReport() {

	bin := r.reporter.falcoBin
	args := r.reporter.falcoargs

	r.reporter.falcoTracer = stats_getter.NewFalcoTracer(r.reporter.mode)

	for i := 0; i < r.iterations; i++ {
		cmd := exec.Command(bin, args...)

		err := cmd.Run()
		if err != nil {
			log.Fatalln("cmd.Run() failed with ", err)
		}

		stats_getter.OpenFalcoGateway(r.reporter.falcoTracer)

		r.getOfflineStats()

		stats_getter.CloseFalcoGateway(r.reporter.falcoTracer)
	}

	r.reporter.falcoTracer.LoadOfflineRulesFromFalco()

	jsonStats, err := r.reporter.falcoTracer.MarshalJSON()
	if err != nil {
		log.Fatal("error in object marshaling")
	}

	writeMetricsOnFile(jsonStats, r.reporter.outputFile)
}

func (r *offlineReporter) getOfflineStats() {
	r.reporter.falcoTracer.LoadOfflineStatsFromFalco()
}
