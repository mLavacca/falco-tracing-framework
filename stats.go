package main

import (
	"strconv"
	"strings"
)

type Falcostats struct {
	StartTime     uint64                 `json:"start time"`
	EndTime       uint64                 `json:"end time"`
	UnbrokenRules map[string]RuleStat    `json:"unbroken rules"`
	BrokenRules   map[string]RuleStat    `json:"broken rules"`
	Stacktraces   map[string]Stacktrace  `json:"stacktraces"`
	CounterStats  map[string]CounterStat `json:"counter statistics"`
}

type Stacktrace struct {
	Counter   uint64     `json:"counter"`
	Functions []FuncStat `json:"functions"`
}

type FuncStat struct {
	Name    string `json:"function"`
	Counter uint64 `json:"counter"`
	Latency uint64 `json:"latency"`
	Level   int    `json:"level"`
}

type CounterStat struct {
	Counter uint64 `json:"counter"`
}

type RuleStat struct {
	Id      int    `json:"rule id"`
	Tag     int    `json:"tag id"`
	Counter uint64 `json:"counter"`
	Latency uint64 `json:"latency"`
}

func NewFalcoStats() *Falcostats {
	f := new(Falcostats)

	f.UnbrokenRules = make(map[string]RuleStat)
	f.BrokenRules = make(map[string]RuleStat)
	f.Stacktraces = make(map[string]Stacktrace)
	f.CounterStats = make(map[string]CounterStat)

	return f
}

func NewStackTrace(line string) (string, *Stacktrace) {
	s := new(Stacktrace)

	line = strings.Replace(line, "\n", "", 1)

	tracerLine := strings.Split(line, "-")

	name := tracerLine[2]
	counter, _ := strconv.ParseUint(tracerLine[3], 10, 64)

	s.Counter = counter

	s.Functions = make([]FuncStat, 0)

	return name, s
}

func NewFuncStat(line string) *FuncStat {
	f := new(FuncStat)

	line = strings.Replace(line, "\n", "", 1)

	tracerLine := strings.Split(line, "-")

	name := tracerLine[1]
	counter, _ := strconv.ParseUint(tracerLine[2], 10, 64)
	latency, _ := strconv.ParseUint(tracerLine[3], 10, 64)
	level, _ := strconv.Atoi(tracerLine[4])

	f.Name = name
	f.Counter = counter
	f.Latency = latency
	f.Level = level

	return f
}

func NewCounterStat(line string) (string, *CounterStat) {
	c := new(CounterStat)

	line = strings.Replace(line, "\n", "", 1)

	tracerLine := strings.Split(line, "-")

	name := tracerLine[1]
	counter, _ := strconv.ParseUint(tracerLine[2], 10, 64)

	c.Counter = counter

	return name, c
}

func NewRuleStat(line string, ra *RulesAggregator) (string, *RuleStat) {
	r := new(RuleStat)

	line = strings.Replace(line, "\n", "", 1)

	tracerLine := strings.Split(line, "-")

	id, _ := strconv.Atoi(tracerLine[1])
	tag, _ := strconv.Atoi(tracerLine[2])
	counter, _ := strconv.ParseUint(tracerLine[3], 10, 64)
	latency, _ := strconv.ParseUint(tracerLine[4], 10, 64)

	rule := ra.getRuleById(id)

	r.Id = rule.Id
	r.Tag = tag
	r.Counter = counter
	r.Latency = latency

	return rule.Name, r
}

func (f *Falcostats) addUnbrokenRule(key string, value RuleStat) {
	f.UnbrokenRules[key] = value
}

func (f *Falcostats) addBrokenRule(key string, value RuleStat) {
	f.BrokenRules[key] = value
}

func (f *Falcostats) addStackTrace(key string, value Stacktrace) {
	f.Stacktraces[key] = value
}

func (f *Falcostats) addFuncStat(key string, value FuncStat) {

	s := f.Stacktraces[key]
	s.Functions = append(f.Stacktraces[key].Functions, value)

	f.Stacktraces[key] = s
}

func (f *Falcostats) addCounterStat(key string, value CounterStat) {
	f.CounterStats[key] = value
}
