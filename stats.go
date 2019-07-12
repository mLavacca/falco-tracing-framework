package main

import (
	"strconv"
	"strings"
)

type Falcostats struct {
	StartTime     uint64
	EndTime       uint64
	UnbrokenRules []RuleStat
	BrokenRules   []RuleStat
	FuncStats     []FuncStat
	CounterStats  []CounterStat
}

type FuncStat struct {
	Name    string
	Counter uint64
	Latency uint64
}

type CounterStat struct {
	Name    string
	Counter uint64
}

type RuleStat struct {
	Name    string
	Tag     int
	Counter uint64
	Latency uint64
}

func NewFuncStat(line string) *FuncStat {
	f := new(FuncStat)

	line = strings.Replace(line, "\n", "", 1)

	tracerLine := strings.Split(line, "-")

	name := tracerLine[1]
	counter, _ := strconv.ParseUint(tracerLine[2], 10, 64)
	latency, _ := strconv.ParseUint(tracerLine[3], 10, 64)

	f.Name = name
	f.Counter = counter
	f.Latency = latency

	return f
}

func NewCounterStat(line string) *CounterStat {
	c := new(CounterStat)

	line = strings.Replace(line, "\n", "", 1)

	tracerLine := strings.Split(line, "-")

	name := tracerLine[1]
	counter, _ := strconv.ParseUint(tracerLine[2], 10, 64)

	c.Name = name
	c.Counter = counter

	return c
}

func NewRuleStat(line string, ra *RulesAggregator) *RuleStat {
	// placeholder, feature to implement
	tag := 0

	r := new(RuleStat)

	line = strings.Replace(line, "\n", "", 1)

	tracerLine := strings.Split(line, "-")

	id, _ := strconv.Atoi(tracerLine[1])
	counter, _ := strconv.ParseUint(tracerLine[2], 10, 64)
	latency, _ := strconv.ParseUint(tracerLine[3], 10, 64)
	name := ra.getRuleNameById(id)

	r.Name = name
	r.Tag = tag
	r.Counter = counter
	r.Latency = latency

	return r
}
