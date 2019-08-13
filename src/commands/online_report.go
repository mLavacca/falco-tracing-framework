package commands

import (
	"configuration"
	"log"
	"os"
	"stats_getter"
	"sync"
)

type onlineReporter struct {
	reporter reporterData
}

func newOnlineReporter(conf configuration.OnlineReportConfiguration) *onlineReporter {
	r := new(onlineReporter)

	r.reporter.falcoBin = conf.ProgConfig.ProgBin
	r.reporter.falcoargs = conf.ProgConfig.ProgArgs
	r.reporter.outputFile = conf.OutputFile
	r.reporter.mode = "online"

	return r
}

func (or *onlineReporter) report() {
	or.startReport()
}

func (r *onlineReporter) startReport() {

	r.reporter.falcoTracer = stats_getter.NewFalcoTracer(r.reporter.mode)

	var wg sync.WaitGroup
	sigs := make(chan os.Signal)

	r.reporter.falcoTracer.LoadRulesFromFalco()

	r.reporter.falcoTracer.FlushFalcoData()

	wg.Add(1)
	//go r.falcoTracer.LoadStatsFromFalco(time.Duration(1), &wg)

	<-sigs

	r.reporter.falcoTracer.ExitFlag = true

	wg.Wait()

	jsonStats, err := r.reporter.falcoTracer.MarshalJSON()
	if err != nil {
		log.Fatal("error in object marshaling")
	}

	r.reporter.falcoTracer.StatsAggregator.SortAvgSlices()

	writeMetricsOnFile(jsonStats, r.reporter.outputFile)
}
