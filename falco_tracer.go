package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type FalcoTracer struct {
	falcoInterface  *FalcoInterface
	statsAggregator *StatsAggregator
	rulesAggregator *RulesAggregator
}

func NewFalcoTracer() *FalcoTracer {

	falcoPid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("Pid invalid")
	}

	f := new(FalcoTracer)

	f.falcoInterface = NewFalcoInterface(falcoPid, "/tmp/tracer_pipe_"+os.Args[1])
	f.statsAggregator = new(StatsAggregator)
	f.rulesAggregator = new(RulesAggregator)

	return f
}

func (f *FalcoTracer) setupConnection() {
	f.falcoInterface.OpenPipe()
}

func (f *FalcoTracer) loadRulesFromFalco() {
	f.falcoInterface.sendSigRcvRulesNames()

	for {
		line := f.falcoInterface.getLine()

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

		f.rulesAggregator.addRule(*r)
	}
}

func (f *FalcoTracer) loadStatsFromFalco(t time.Duration, ch chan (StatsAggregator)) {
	for {
		f.falcoInterface.sendSigRcvSummary()

		for {
			line := f.falcoInterface.getLine()

			if strings.Contains(string(line), "START SUMMARY") {
				start, end := getTimesFromMessage(line)
				f.statsAggregator.addFalcoStats(start, end)
				continue
			}

			if strings.Contains(string(line), "START FUNCTIONS LATENCIES") {
				f.getFunctionLatencies()
				continue
			}

			if strings.Contains(string(line), "START COUNTERS") {
				f.getCounters()
				continue
			}

			if strings.Contains(string(line), "START UNBROKEN RULES") {
				f.getUnbrokenRules()
				continue
			}

			if strings.Contains(string(line), "START BROKEN RULES") {
				f.getBrokenRules()
				continue
			}

			if strings.Contains(string(line), "END SUMMARY") {
				f.falcoInterface.sendSigFlushData()
				break
			}
		}

		ch <- *f.statsAggregator
		time.Sleep(t * time.Second)
	}
}

func (f *FalcoTracer) getFunctionLatencies() {
	for {
		line := f.falcoInterface.getLine()

		if strings.Contains(string(line), "END FUNCTIONS LATENCIES") {
			break
		}

		fs := NewFuncStat(line)

		f.statsAggregator.addFuncStat(*fs)
	}
}

func (f *FalcoTracer) getCounters() {
	for {
		line := f.falcoInterface.getLine()

		if strings.Contains(string(line), "END COUNTERS") {
			break
		}

		cs := NewCounterStat(line)
		f.statsAggregator.addCounterStat(*cs)
	}
}

func (f *FalcoTracer) getUnbrokenRules() {

	for {
		line := f.falcoInterface.getLine()

		if strings.Contains(string(line), "END UNBROKEN RULES") {
			break
		}

		ur := NewRuleStat(line, f.rulesAggregator)
		f.statsAggregator.addUnbrokenRuleStat(*ur)
	}
}

func (f *FalcoTracer) getBrokenRules() {

	for {
		line := f.falcoInterface.getLine()

		if strings.Contains(string(line), "END BROKEN RULES") {
			break
		}

		br := NewRuleStat(line, f.rulesAggregator)
		f.statsAggregator.addBrokenRuleStat(*br)
	}
}

func getTimesFromMessage(line string) (uint64, uint64) {

	line = strings.Replace(line, "\n", "", 1)

	tracerLine := strings.Split(line, "-")

	start, _ := strconv.ParseUint(tracerLine[2], 10, 64)
	end, _ := strconv.ParseUint(tracerLine[3], 10, 64)

	return start, end
}
