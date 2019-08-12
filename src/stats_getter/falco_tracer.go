package stats_getter

import (
	"encoding/json"
	"strconv"
	"strings"

	m "metrics"
)

type FalcoTracer struct {
	ExitFlag        bool
	falcoGateway    *FalcoGateway
	StatsAggregator *m.StatsAggregator
	rulesAggregator *m.RulesAggregator
}

func NewFalcoTracer(mode string) *FalcoTracer {

	f := new(FalcoTracer)

	f.falcoGateway = NewFalcoGateway(mode)
	f.StatsAggregator = new(m.StatsAggregator)
	f.rulesAggregator = m.NewRulesAggregator()

	return f
}

func (f *FalcoTracer) LoadRulesFromFalco() {
	f.falcoGateway.sendSigRcvRulesNames()

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

		rs := m.NewruleStatavg(r.Name, r.Id)
		f.StatsAggregator.AvgUnbrokenRulesStats = append(f.StatsAggregator.AvgUnbrokenRulesStats, *rs)
	}

	f.rulesAggregator.SetNRules()
}

func (f *FalcoTracer) FlushFalcoData() {
	f.falcoGateway.sendSigFlushData()
}

func (f *FalcoTracer) LoadStatsFromFalco( /*t time.Duration, wg *sync.WaitGroup*/ ) {
	for f.ExitFlag == false {
		//time.Sleep(t * time.Second)

		//f.falcoGateway.sendSigRcvSummary()

		for {
			line := f.falcoGateway.getLine()

			if strings.Contains(string(line), "START SUMMARY") {
				start, end := getTimesFromMessage(line)
				f.StatsAggregator.AddFalcoStats(start, end)
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
				//f.FlushFalcoData()
				if f.falcoGateway.mode == "offline" {
					return
				}
				break
			}
		}
	}

	f.StatsAggregator.SetTimes()

	//wg.Done()
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
				f.StatsAggregator.AddStackTrace(name, *s)
			} else {
				nameFunc, fs := m.NewFuncStat(line)
				f.StatsAggregator.AddFuncStat(name, nameFunc, *fs)
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

		f.StatsAggregator.AddCounterStat(name, *cs)
	}
}

func (f *FalcoTracer) getUnbrokenRules() {

	for {
		line := f.falcoGateway.getLine()

		if strings.Contains(string(line), "END UNBROKEN RULES") {
			break
		}

		name, ur := m.NewRuleStat(line, f.rulesAggregator)
		f.StatsAggregator.AddUnbrokenRuleStat(name, *ur)

		f.StatsAggregator.SumValuesToAverageUnbroken(ur.Id, ur.Counter, ur.Latency)
	}
}

func (f *FalcoTracer) getBrokenRules() {

	for {
		line := f.falcoGateway.getLine()

		if strings.Contains(string(line), "END BROKEN RULES") {
			break
		}

		name, br := m.NewRuleStat(line, f.rulesAggregator)
		f.StatsAggregator.AddBrokenRuleStat(name, *br)

		f.StatsAggregator.SumValuesToAverageBroken(br.Id, br.Counter, br.Latency)
	}
}

func getTimesFromMessage(line string) (uint64, uint64) {

	line = strings.Replace(line, "\n", "", 1)

	tracerLine := strings.Split(line, "-")

	start, _ := strconv.ParseUint(tracerLine[2], 10, 64)
	end, _ := strconv.ParseUint(tracerLine[3], 10, 64)

	return start, end
}

func (f *FalcoTracer) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Stats m.StatsAggregator `json:"statistics"`
		Rules m.RulesAggregator `json:"rules"`
	}{
		Stats: *f.StatsAggregator,
		Rules: *f.rulesAggregator,
	})
}
