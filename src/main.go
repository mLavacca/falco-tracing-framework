package main

import (
	"commands"
	"configuration"
	"flag"
	"log"
)

func main() {

	var recorder *commands.Recorder
	var reporter *commands.Reporter

	var configFile = flag.String("c", "config.yaml", "configuration yaml file")
	var mode = flag.String("m", "record", "working mode: record, report or compare")

	flag.Parse()

	tracerConf := new(configuration.TracerConfigurations)
	err := tracerConf.UnmarshalYAML(*configFile)
	if err != nil {
		log.Fatalln("unable to load tracer configuration", err)
	}

	switch *mode {
	case "record":
		recorder = commands.NewRecorder(*tracerConf)
		recorder.StartRecord()
		recorder.Rollback()
	case "report":
		reporter = commands.NewReporter(*tracerConf)
		reporter.StartReport()
	}
}
