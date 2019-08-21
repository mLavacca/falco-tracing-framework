package main

import (
	cmd "commands"
	"configuration"
	"flag"
	"log"
)

func main() {

	var configFile = flag.String("c", "config.yaml", "configuration yaml file")
	var command = flag.String("m", "record", "working mode: record, report or compare")

	flag.Parse()

	tracerConf := new(configuration.TracerConfigurations)
	err := tracerConf.UnmarshalYAML(*configFile)
	if err != nil {
		log.Fatalln("unable to load tracer configuration", err)
	}

	cmd.DispatchCommand(*command, *tracerConf)
}
