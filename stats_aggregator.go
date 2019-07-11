package main

import (
	"strconv"
	"strings"
)

type StatsAggregator struct {
	falcoStats Falcostats
}

func (s *StatsAggregator) addUnbrokenRuleStat(name string, tag int, counter uint64, latency uint64) {
	r := NewRuleStat(name, tag, counter, latency)

	s.falcoStats.unbrokenRules = append(s.falcoStats.unbrokenRules, *r)
}

func (s *StatsAggregator) addBrokenRuleStat(name string, tag int, counter uint64, latency uint64) {
	r := NewRuleStat(name, tag, counter, latency)

	s.falcoStats.brokenRules = append(s.falcoStats.unbrokenRules, *r)
}

func (s *StatsAggregator) addFuncStat(name string, counter uint64, latency uint64) {
	f := NewFuncStat(name, counter, latency)

	s.falcoStats.funcStats = append(s.falcoStats.funcStats, *f)
}

func (s *StatsAggregator) addCounterStat(name string, tag int, counter uint64, latency uint64) {
	c := NewCounterStat(name, counter)

	s.falcoStats.counterStats = append(s.falcoStats.counterStats, *c)
}

func (s *StatsAggregator) getFunctionLatencies(falcoInterface *FalcoInterface) {
	for {
		line := falcoInterface.getLine()

		if strings.Contains(string(line), "TRACER INFO - END FUNCTIONS LATENCIES") {
			break
		}

		line = strings.Replace(line, "\n", "", 1)

		tracerLine := strings.Split(line, "-")

		name := tracerLine[1]
		counter, _ := strconv.ParseUint(tracerLine[2], 10, 64)
		latency, _ := strconv.ParseUint(tracerLine[3], 10, 64)

		var f FuncStat
		f.name = name
		f.counter = counter
		f.latency = latency

		s.falcoStats.funcStats = append(s.falcoStats.funcStats, f)
	}
}

func (s *StatsAggregator) getCounters(falcoInterface *FalcoInterface) {
	for {
		line := falcoInterface.getLine()

		if strings.Contains(string(line), "TRACER INFO - END COUNTERS") {
			break
		}

		line = strings.Replace(line, "\n", "", 1)

		tracerLine := strings.Split(line, "-")

		name := tracerLine[1]
		counter, _ := strconv.ParseUint(tracerLine[2], 10, 64)

		var c CounterStat
		c.name = name
		c.counter = counter

		s.falcoStats.counterStats = append(s.falcoStats.counterStats, c)
	}
}

func (s *StatsAggregator) getUnbrokenRules(falcoInterface *FalcoInterface, falcoRules []FalcoRule) {

	for {
		line := falcoInterface.getLine()

		if strings.Contains(string(line), "TRACER INFO - END UNBROKEN RULES") {
			break
		}

		line = strings.Replace(line, "\n", "", 1)

		tracerLine := strings.Split(line, "-")

		id, _ := strconv.Atoi(tracerLine[2])
		counter, _ := strconv.ParseUint(tracerLine[2], 10, 64)
		latency, _ := strconv.ParseUint(tracerLine[3], 10, 64)

		var r RuleStat

		for i := 0; i < len(falcoRules); i++ {
			if falcoRules[i].id == id {
				r.name = falcoRules[i].name
				break
			}
		}

		r.counter = counter
		r.latency = latency

		s.falcoStats.unbrokenRules = append(s.falcoStats.unbrokenRules, r)
	}
}

func (s *StatsAggregator) getBrokenRules(falcoInterface *FalcoInterface, falcoRules []FalcoRule) {

	for {
		line := falcoInterface.getLine()

		if strings.Contains(string(line), "TRACER INFO - END BROKEN RULES") {
			break
		}

		line = strings.Replace(line, "\n", "", 1)

		tracerLine := strings.Split(line, "-")

		id, _ := strconv.Atoi(tracerLine[2])
		counter, _ := strconv.ParseUint(tracerLine[2], 10, 64)
		latency, _ := strconv.ParseUint(tracerLine[3], 10, 64)

		var r RuleStat

		for i := 0; i < len(falcoRules); i++ {
			if falcoRules[i].id == id {
				r.name = falcoRules[i].name
				break
			}
		}

		r.counter = counter
		r.latency = latency

		s.falcoStats.brokenRules = append(s.falcoStats.brokenRules, r)
	}
}
