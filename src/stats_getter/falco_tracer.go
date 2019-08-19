package stats_getter

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	m "metrics"
)

type FalcoTracer struct {
	ExitFlag        bool
	falcoGateway    *FalcoGateway
	offlineMetrics  *m.OfflineMetrics
	rulesAggregator *m.RulesAggregator
}

func NewFalcoTracer(mode string) *FalcoTracer {

	f := new(FalcoTracer)

	f.falcoGateway = newFalcoGateway(mode)

	if f.falcoGateway.mode == "online" {
	}

	if f.falcoGateway.mode == "offline" {
		f.offlineMetrics = new(m.OfflineMetrics)
	}

	f.rulesAggregator = m.NewRulesAggregator()

	return f
}

func OpenFalcoGateway(f *FalcoTracer) {
	f.falcoGateway.openPipe()
}

func CloseFalcoGateway(f *FalcoTracer) {
	f.falcoGateway.closePipe()
}

func (f *FalcoTracer) LoadOfflineRulesFromFalco() {
	f.falcoGateway.openPipeForRules()
	f.loadRulesFromFalco()
	f.falcoGateway.closePipe()
}

func (f *FalcoTracer) LoadOnlineRulesFromFalco() {
	f.falcoGateway.sendSigRcvRulesNames()
	f.loadRulesFromFalco()
}

func (f *FalcoTracer) loadRulesFromFalco() {

	for {
		line := f.falcoGateway.getLine()

		if strings.Contains(line, "START RULES NAMES") {
			continue
		}

		if strings.Contains(line, "END RULES NAMES") {
			break
		}

		r := m.NewRule(line)
		if r == nil {
			continue
		}

		f.rulesAggregator.AddRule(*r)
	}

	f.rulesAggregator.SetNRules()
}

func (f *FalcoTracer) FlushFalcoData() {
	f.falcoGateway.sendSigFlushData()
}

func (f *FalcoTracer) loadStatsFromFalco() {
	for {
		line := f.falcoGateway.getLine()

		if strings.Contains(string(line), "START SUMMARY") {
			//start, end := getTimesFromMessage(line)
			f.offlineMetrics.AddOfflineMetrics()
			continue
		}

		if strings.Contains(string(line), "START STACKTRACES") {
			f.getStacktraces()
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
			break
		}
	}
}

func (f *FalcoTracer) LoadOfflineStatsFromFalco() {
	f.loadStatsFromFalco()
}

func (f *FalcoTracer) LoadOnlineStatsFromFalco(t time.Duration, wg *sync.WaitGroup) {
	for f.ExitFlag == false {
		time.Sleep(t * time.Second)

		f.loadStatsFromFalco()

		f.FlushFalcoData()

		f.falcoGateway.sendSigRcvSummary()
	}

	wg.Done()
}

func (f *FalcoTracer) getStacktraces() {
	for {
		line := f.falcoGateway.getLine()

		if strings.Contains(string(line), "END STACKTRACES") {
			break
		}

		var name string
		var s *m.Stacktrace

		for {
			if strings.Contains(string(line), "END STACKTRACE") {
				break
			}

			if strings.Contains(string(line), "START STACKTRACE") {
				name, s = m.NewStackTrace(line)

				if f.falcoGateway.mode == "offline" {
					f.offlineMetrics.AddStackTrace(name, *s)
				}

			} else {
				nameFunc, fs := m.NewFuncStat(line)

				if f.falcoGateway.mode == "offline" {
					f.offlineMetrics.AddFuncStat(name, nameFunc, *fs)
				}
			}

			line = f.falcoGateway.getLine()
		}
	}
}

func (f *FalcoTracer) getCounters() {
	for {
		line := f.falcoGateway.getLine()

		if strings.Contains(string(line), "END COUNTERS") {
			break
		}

		name, cs := m.NewCounterStat(line)

		if f.falcoGateway.mode == "offline" {
			f.offlineMetrics.AddCounterStat(name, *cs)
		}
	}
}

func (f *FalcoTracer) getUnbrokenRules() {

	for {
		line := f.falcoGateway.getLine()

		if strings.Contains(string(line), "END UNBROKEN RULES") {
			break
		}

		name, ur := m.NewRuleStat(line, f.rulesAggregator)

		if f.falcoGateway.mode == "online" {
		}

		if f.falcoGateway.mode == "offline" {
			f.offlineMetrics.AddUnbrokenRuleMetric(name, *ur)
		}
	}
}

func (f *FalcoTracer) getBrokenRules() {

	for {
		line := f.falcoGateway.getLine()

		if strings.Contains(string(line), "END BROKEN RULES") {
			break
		}

		name, ur := m.NewRuleStat(line, f.rulesAggregator)

		if f.falcoGateway.mode == "online" {
		}

		if f.falcoGateway.mode == "offline" {
			f.offlineMetrics.AddBrokenRuleMetric(name, *ur)
		}
	}
}

func getTimesFromMessage(line string) (uint64, uint64) {

	line = strings.Replace(line, "\n", "", 1)

	tracerLine := strings.Split(line, "-")

	start, _ := strconv.ParseUint(tracerLine[2], 10, 64)
	end, _ := strconv.ParseUint(tracerLine[3], 10, 64)

	return start, end
}

func (f *FalcoTracer) OfflineAvg() m.OfflineFalcoMetrics {
	var metr m.FalcoMetrics = f.offlineMetrics.Fm[0].Metrics
	var avgUnbroken = f.offlineMetrics.Fm[0].UnbrokenRuleMetrics
	var avgBroken = f.offlineMetrics.Fm[0].BrokenRuleMetrics

	// stacktrace avg computation
	for i, om := range f.offlineMetrics.Fm[1:] {
		sts1 := metr.Stacktraces
		sts2 := om.Metrics.Stacktraces

		var j uint64 = uint64(i) + 1

		stacktraceMap := make(map[string]m.Stacktrace)
		for k, v := range sts2 {

			var st m.Stacktrace
			st.Counter = v.Counter

			funcMap := make(map[string]m.FuncStat)
			for ki, vi := range v.Functions {
				var f m.FuncStat

				f.Latency = (vi.Latency + (sts1[k].Functions[ki].Latency * j)) / (j + 1)
				f.Caller = vi.Caller

				funcMap[ki] = f
			}

			st.Functions = funcMap
			stacktraceMap[k] = st
		}

		metr.Stacktraces = stacktraceMap
	}

	// counter avg computation
	for i, om := range f.offlineMetrics.Fm[1:] {
		cts1 := metr.CounterStats
		cts2 := om.Metrics.CounterStats

		var j uint64 = uint64(i) + 1

		counterMap := make(map[string]m.CounterStat)
		for k, v := range cts2 {
			var c m.CounterStat

			c.Counter = (v.Counter + (cts1[k].Counter * j)) / (j + 1)
			counterMap[k] = c
		}
		metr.CounterStats = counterMap
	}

	// unbroken rules avg computation
	for i, om := range f.offlineMetrics.Fm[1:] {
		urm1 := avgUnbroken
		urm2 := om.UnbrokenRuleMetrics

		var j uint64 = uint64(i) + 1

		ruleMap := make(map[string]m.RuleStat)
		for k, v := range urm2 {
			var r m.RuleStat

			r.Id = v.Id
			r.Tag = v.Tag
			r.Counter = (v.Counter + (urm1[k].Counter * j)) / (j + 1)
			r.Latency = (v.Latency + (urm1[k].Latency * j)) / (j + 1)
			ruleMap[k] = r
		}
		avgUnbroken = ruleMap
	}

	// broken rules avg computation
	for i, om := range f.offlineMetrics.Fm[1:] {
		brm1 := avgBroken
		brm2 := om.BrokenRuleMetrics

		var j uint64 = uint64(i) + 1

		ruleMap := make(map[string]m.RuleStat)
		for k, v := range brm2 {
			var r m.RuleStat

			r.Id = v.Id
			r.Tag = v.Tag
			r.Counter = (v.Counter + (brm1[k].Counter * j)) / (j + 1)
			r.Latency = (v.Latency + (brm1[k].Latency * j)) / (j + 1)

			ruleMap[k] = r
		}
		avgBroken = ruleMap
	}

	return m.OfflineFalcoMetrics{
		Metrics:             metr,
		UnbrokenRuleMetrics: avgUnbroken,
		BrokenRuleMetrics:   avgBroken,
	}
}

func (f *FalcoTracer) MarshalOfflineJSON(metr m.OfflineFalcoMetrics) ([]byte, error) {
	return json.Marshal(metr)
}
