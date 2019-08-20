package commands

import (
	"configuration"
	"log"
	"os/exec"
	"stats_getter"

	df "data_formatter"
)

type offlineReporter struct {
	reporter reporterData

	iterations int
}

func newOfflineReporter(conf configuration.OfflineReportConfiguration) *offlineReporter {
	r := new(offlineReporter)

	r.reporter.falcoBins = conf.ProgConfig.ProgBins
	r.reporter.falcoargs = conf.ProgConfig.ProgArgs
	r.reporter.outputFile = conf.OutputFile
	r.reporter.outputFoldedFile = conf.OutputFoldedStacktrace
	r.reporter.outputDottedFile = conf.OutputDottedStacktrace

	r.reporter.mode = "offline"

	r.iterations = conf.Iterations

	return r
}

func (or *offlineReporter) report() {
	or.startReport()
}

func (r *offlineReporter) startReport() {

	r.reporter.falcoTracer = stats_getter.NewFalcoTracer(r.reporter.mode)

	for _, bin := range r.reporter.falcoBins {
		args := r.reporter.falcoargs

		for i := 0; i < r.iterations; i++ {
			cmd := exec.Command(bin, args...)

			err := cmd.Run()
			if err != nil {
				log.Fatalln("cmd.Run() failed with ", err)
			}

			stats_getter.OpenFalcoGateway(r.reporter.falcoTracer)

			r.reporter.falcoTracer.LoadOfflineRulesFromFalco()

			stats_getter.OpenFalcoGateway(r.reporter.falcoTracer)

			r.getOfflineStats()

			stats_getter.CloseFalcoGateway(r.reporter.falcoTracer)
		}
	}

	metr := r.reporter.falcoTracer.OfflineAvg()

	jsonStats, err := r.reporter.falcoTracer.MarshalOfflineJSON(metr)
	if err != nil {
		log.Fatal("error in object marshaling")
	}

	foldedStacktraces := df.CreateFoldedStacktrace(metr.Metrics.Stacktraces)
	dottedStackTrace := df.CreateDotStacktrace(metr.Metrics.Stacktraces)

	writeMetricsOnFile(dottedStackTrace, r.reporter.outputDottedFile)
	writeMetricsOnFile(jsonStats, r.reporter.outputFile)
	writeMetricsOnFile(foldedStacktraces, r.reporter.outputFoldedFile)

}

func (r *offlineReporter) getOfflineStats() {
	r.reporter.falcoTracer.LoadOfflineStatsFromFalco()
}
