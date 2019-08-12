package commands

import (
	"configuration"
	"fmt"
	"log"
	"os"
	"os/exec"
	"stats_getter"
)

type Reporter struct {
	falcoBin  string
	falcoargs []string

	mode       string
	iterations int

	falcoTracer *stats_getter.FalcoTracer
}

func NewReporter(conf configuration.TracerConfigurations) *Reporter {
	r := new(Reporter)

	r.falcoBin = conf.Report.ProgConfig.ProgBin
	r.falcoargs = conf.Report.ProgConfig.ProgArgs
	r.mode = conf.Report.Mode
	r.iterations = conf.Report.Iterations

	//r.falcoTracer = stats_getter.NewFalcoTracer(r.mode)

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

func (r *Reporter) offlineReport() {
	cmd := exec.Command(r.falcoBin, r.falcoargs...)

	/*
		for i := 0; i < r.iterations; i++ {
			err := cmd.Run()
			if err != nil {
				log.Fatalln("cmd.Start() failed with ", err)
			}

			getOfflineStats()
		}
	*/

	err := cmd.Start()
	if err != nil {
		log.Fatalln("cmd.Start() failed with ", err)
	}

	r.falcoTracer = stats_getter.NewFalcoTracer(r.mode)
	r.getOfflineStats()
}

func (r *Reporter) onlineReport() {
	/*
		var wg sync.WaitGroup
		sigs := make(chan os.Signal)

		r.falcoTracer.LoadRulesFromFalco()

		r.falcoTracer.FlushFalcoData()

		wg.Add(1)
		go r.falcoTracer.LoadStatsFromFalco(time.Duration(1), &wg)

		<-sigs

		r.falcoTracer.ExitFlag = true

		wg.Wait()

		jsonStats, err := r.falcoTracer.MarshalJSON()
		if err != nil {
			log.Fatal("error in object marshaling")
		}

		r.falcoTracer.StatsAggregator.SortAvgSlices()

		writeJSONOnFile(jsonStats, "/tmp/")
	*/
}

func (r *Reporter) getOfflineStats() {
	//r.falcoTracer.LoadRulesFromFalco()
	r.falcoTracer.LoadStatsFromFalco()

	jsonStats, err := r.falcoTracer.MarshalJSON()
	if err != nil {
		log.Fatal("error in object marshaling")
	}

	r.falcoTracer.StatsAggregator.SortAvgSlices()

	writeJSONOnFile(jsonStats, "/tmp/")

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
