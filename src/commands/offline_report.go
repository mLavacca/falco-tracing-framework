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

	cmd := exec.Command(bin, args...)

	r.iterations = 1

	for i := 0; i < r.iterations; i++ {
		err := cmd.Run()
		if err != nil {
			log.Fatalln("cmd.Start() failed with ", err)
		}

		r.reporter.falcoTracer = stats_getter.NewFalcoTracer(r.reporter.mode)

		r.getOfflineStats()

		stats_getter.CloseFalcoTracer(r.reporter.falcoTracer)
	}

}

func (r *offlineReporter) getOfflineStats() {
	r.reporter.falcoTracer.LoadOfflineStatsFromFalco()

	jsonStats, err := r.reporter.falcoTracer.MarshalJSON()
	if err != nil {
		log.Fatal("error in object marshaling")
	}

	r.reporter.falcoTracer.StatsAggregator.SortAvgSlices()

	writeMetricsOnFile(jsonStats, r.reporter.outputFile)
}
