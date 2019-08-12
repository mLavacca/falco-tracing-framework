package metrics

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
	Counter   uint64              `json:"counter"`
	Functions map[string]FuncStat `json:"functions"`
}

type FuncStat struct {
	Counter uint64 `json:"counter"`
	Latency uint64 `json:"latency"`
	Caller  string `json:"caller"`
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

type RuleStatAvg struct {
	Rule RuleStat `json:"values"`
	Name string   `json:"name"`
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

	s.Functions = make(map[string]FuncStat)

	return name, s
}

func NewFuncStat(line string) (string, *FuncStat) {
	f := new(FuncStat)

	line = strings.Replace(line, "\n", "", 1)

	tracerLine := strings.Split(line, "-")

	name := tracerLine[1]
	counter, _ := strconv.ParseUint(tracerLine[2], 10, 64)
	latency, _ := strconv.ParseUint(tracerLine[3], 10, 64)
	caller := tracerLine[4]

	f.Counter = counter
	f.Latency = latency
	f.Caller = caller

	return name, f
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

func NewruleStatavg(name string, id int) *RuleStatAvg {
	rs := new(RuleStatAvg)
	r := new(RuleStat)

	rs.Rule = *r
	rs.Name = name
	rs.Rule.Id = id

	return rs
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

func (f *Falcostats) addFuncStat(key string, funcName string, value FuncStat) {
	f.Stacktraces[key].Functions[funcName] = value
}

func (f *Falcostats) addCounterStat(key string, value CounterStat) {
	f.CounterStats[key] = value
}
