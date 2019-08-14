package metrics

import (
	"strconv"
	"strings"
)

type FalcoMetrics struct {
	Stacktraces  map[string]Stacktrace  `json:"stacktraces"`
	CounterStats map[string]CounterStat `json:"counter statistics"`
}

type FalcoRulesMetrics struct {
	UnbrokenRules map[string]RuleStat `json:"unbroken rules"`
	BrokenRules   map[string]RuleStat `json:"broken rules"`
}

type OnlineFalcoMetrics struct {
	StartTime uint64 `json:"start time"`
	EndTime   uint64 `json:"end time"`

	Metrics FalcoMetrics `json:"falco metrics"`
}

type OfflineFalcoMetrics struct {
	Metrics FalcoMetrics `json:"falco metrics"`
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

func NewOfflineFalcoMetrics() *OfflineFalcoMetrics {
	m := new(OfflineFalcoMetrics)

	m.Metrics.Stacktraces = make(map[string]Stacktrace)
	m.Metrics.CounterStats = make(map[string]CounterStat)

	return m
}

func NewOnlineFalcoMetrics() *OnlineFalcoMetrics {
	m := new(OnlineFalcoMetrics)

	m.Metrics.Stacktraces = make(map[string]Stacktrace)
	m.Metrics.CounterStats = make(map[string]CounterStat)

	return m
}

func getTracerLine(line string) []string {
	line = strings.Replace(line, "\n", "", 1)

	tracerLine := strings.Split(line, "-")
	return tracerLine
}

func NewStackTrace(line string) (string, *Stacktrace) {
	s := new(Stacktrace)

	tracerLine := getTracerLine(line)

	name := tracerLine[2]
	counter, _ := strconv.ParseUint(tracerLine[3], 10, 64)

	s.Counter = counter

	s.Functions = make(map[string]FuncStat)

	return name, s
}

func NewFuncStat(line string) (string, *FuncStat) {
	f := new(FuncStat)

	tracerLine := getTracerLine(line)

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

	tracerLine := getTracerLine(line)

	name := tracerLine[1]
	counter, _ := strconv.ParseUint(tracerLine[2], 10, 64)

	c.Counter = counter

	return name, c
}

func NewRuleStat(line string, ra *RulesAggregator) (string, *RuleStat) {
	r := new(RuleStat)

	tracerLine := getTracerLine(line)

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

func (f *OfflineFalcoMetrics) addStackTrace(key string, value Stacktrace) {
	f.Metrics.Stacktraces[key] = value
}

func (f *OfflineFalcoMetrics) addFuncStat(key string, funcName string, value FuncStat) {
	f.Metrics.Stacktraces[key].Functions[funcName] = value
}

func (f *OfflineFalcoMetrics) addCounterStat(key string, value CounterStat) {
	f.Metrics.CounterStats[key] = value
}
