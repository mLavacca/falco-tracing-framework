package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	var falcoInterface *FalcoInterface
	var falcoRules []FalcoRule
	var statsAggregator *StatsAggregator = new(StatsAggregator)

	falcoPid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("Pid invalid")
	}

	falcoInterface = NewFalcoInterface(falcoPid, "/tmp/tracer_pipe_"+os.Args[1])

	falcoInterface.OpenPipe()

	falcoInterface.sendSigRcvRulesNames()

	for {
		line := falcoInterface.getLine()

		if strings.Contains(line, "TRACER INFO - START RULES NAMES") {
			continue
		}

		if strings.Contains(line, "TRACER INFO - END RULES NAMES") {
			break
		}

		r := NewRule(line)
		if r == nil {
			continue
		}

		falcoRules = append(falcoRules, *r)
	}

	for {
		falcoInterface.sendSigRcvSummary()

		for {
			line := falcoInterface.getLine()

			if strings.Contains(string(line), "TRACER INFO - START SUMMARY") {
				continue
			}

			if strings.Contains(string(line), "TRACER INFO - START FUNCTIONS LATENCIES") {
				statsAggregator.getFunctionLatencies(falcoInterface)
				continue
			}

			if strings.Contains(string(line), "TRACER INFO - START COUNTERS") {
				statsAggregator.getCounters(falcoInterface)
				continue
			}

			if strings.Contains(string(line), "TRACER INFO - START UNBROKEN RULES") {
				statsAggregator.getUnbrokenRules(falcoInterface, falcoRules)
				continue
			}

			if strings.Contains(string(line), "TRACER INFO - START BROKEN RULES") {
				statsAggregator.getBrokenRules(falcoInterface, falcoRules)
				continue
			}

			if strings.Contains(string(line), "TRACER INFO - END SUMMARY") {
				break
			}
		}

		time.Sleep(4 * time.Second)
	}
}
