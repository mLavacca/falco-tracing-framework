package commands

import (
	"configuration"
	"fmt"
	"log"
	"os"
	"os/exec"
	"stats_getter"
	"time"
)

type Reporter struct {
	falcoBin  string
	falcoargs []string

	mode       string
	iterations int

	falcoTracer *stats_getter.falcoTracer
}

func NewReporter(conf configuration.TracerConfigurations) *Reporter {
	r := new(Reporter)

	r.falcoBin = conf.Report.ProgConfig.ProgBin
	r.falcoargs = conf.Report.ProgConfig.ProgArgs
	r.mode = conf.Report.Mode
	r.iterations = conf.Report.Iterations
	r.falcoTracer = stats_getter.NewFalcoTracer()

	return r
}

func (r *Reporter) StartReport() {

	if r.mode == "offline" {
		r.offlineReport()
	}

	if r.mode == "online" {
		r.onlineReport()
	}
}

func (r *reporter) offlineReport() {
	cmd := exec.Command(r.falcoBin, r.falcoargs...)

	for i := 0; i < r.iterations; i++ {
		err := cmd.Run()
		if err != nil {
			log.Fatalln("cmd.Start() failed with ", err)
		}

		getOfflineStats()
	}
}

func (r *reporter) onlineReport() {
	falcoTracer.setupConnection()

	falcoTracer.loadRulesFromFalco()

	falcoTracer.flushFalcoData()

	wg.Add(1)
	go falcoTracer.loadStatsFromFalco(time.Duration(t), &wg)

	<-sigs

	falcoTracer.exitFlag = true

	wg.Wait()

	jsonStats, err := falcoTracer.MarshalJSON()
	if err != nil {
		log.Fatal("error in object marshaling")
	}

	falcoTracer.statsAggregator.sortAvgSlices()

	writeJSONOnFile(jsonStats, outDir)
}

func getOfflineStats() {

}

func writeJSONOnFile(jsonStats []byte, outDir string) {
	f, err := os.Create(outDir + "tracer_data.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	l, err := f.Write(jsonStats)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
