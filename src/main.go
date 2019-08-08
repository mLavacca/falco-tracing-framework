package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func main() {

	var recorder *Recorder
	var reporter *Reporter

	var configFile = flag.String("c", "config.yaml", "configuration yaml file")
	var mode = flag.String("m", "record", "working mode")
	var outputFile = flag.String("o", "", "output file")

	flag.Parse()

	tracerConf := new(TracerConfigurations)
	err := tracerConf.UnmarshalYAML(*configFile)
	if err != nil {
		log.Fatalln("unable to load Falco configuration", err)
	}

	switch *mode {
	case "record":
		if *outputFile == "" {
			*outputFile = "./trace.scap"
		}
		recorder = NewRecorder(*tracerConf)
		recorder.startRecord()
	case "report":
		if *outputFile == "" {
			*outputFile = "./stats.json"
		}
		reporter = NewReporter(tracerConf.Report)
		reporter.startReport()
	}
}

func tmpThrash() {

	var wg sync.WaitGroup

	outDir := "/tmp/"

	t, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal("error, time parameter missing")
	}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	falcoTracer := NewFalcoTracer()

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

	dumpJSONOnFile(jsonStats, outDir)
}

func dumpJSONOnFile(jsonStats []byte, outDir string) {
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
