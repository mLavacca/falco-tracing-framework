package stats_getter

import (
	"encoding/json"
	"log"
	"metrics"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type FalcoTracer struct {
	exitFlag        bool
	falcoGateway    *FalcoGateway
	statsAggregator *metrics.StatsAggregator
	rulesAggregator *metrics.RulesAggregator
}

func NewFalcoTracer() *FalcoTracer {

	falcoPid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("Pid invalid")
	}

	f := new(FalcoTracer)

	f.falcoGateway = NewFalcoGateway(falcoPid, "/tmp/tracer_pipe_"+os.Args[1])
	f.statsAggregator = new(StatsAggregator)
	f.rulesAggregator = NewRulesAggregator()

	return f
}

func (f *FalcoTracer) setupConnection() {
	f.falcoGateway.OpenPipe()
}

func (f *FalcoTracer) loadRulesFromFalco() {
	f.falcoGateway.sendSigRcvRulesNames()

	for {
		line := f.falcoGateway.getLine()

		if strings.Contains(line, "START RULES NAMES") {
			continue
		}

		if strings.Contains(line, "END RULES NAMES") {
			break
		}

		r := NewRule(line)
		if r == nil {
			continue
		}

		f.rulesAggregator.addRule(*r)

		rs := NewruleStatavg(r.Name, r.Id)
		f.statsAggregator.AvgUnbrokenRulesStats = append(f.statsAggregator.AvgUnbrokenRulesStats, *rs)
	}

	f.rulesAggregator.setNRules()
}

func (f *FalcoTracer) flushFalcoData() {
	f.falcoGateway.sendSigFlushData()
}

func (f *FalcoTracer) loadStatsFromFalco(t time.Duration, wg *sync.WaitGroup) {
	for f.exitFlag == false {
		time.Sleep(t * time.Second)

		f.falcoGateway.sendSigRcvSummary()

		for {
			line := f.falcoGateway.getLine()

			if strings.Contains(string(line), "START SUMMARY") {
				start, end := getTimesFromMessage(line)
				f.statsAggregator.addFalcoStats(start, end)
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
				f.flushFalcoData()
				break
			}
		}
	}

	f.statsAggregator.setTimes()

	wg.Done()
}

func (f *FalcoTracer) getStacktraces() {
	for {
		line := f.falcoGateway.getLine()

		if strings.Contains(string(line), "END STACKTRACES") {
			break
		}

		var name string
		var s *Stacktrace

		for {
			if strings.Contains(string(line), "END STACKTRACE") {
				break
			}

			if strings.Contains(string(line), "START STACKTRACE") {
				name, s = NewStackTrace(line)
				f.statsAggregator.addStackTrace(name, *s)
			} else {
				nameFunc, fs := NewFuncStat(line)
				f.statsAggregator.addFuncStat(name, nameFunc, *fs)
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

		name, cs := NewCounterStat(line)

		f.statsAggregator.addCounterStat(name, *cs)
	}
}

func (f *FalcoTracer) getUnbrokenRules() {

	for {
		line := f.falcoGateway.getLine()

		if strings.Contains(string(line), "END UNBROKEN RULES") {
			break
		}

		name, ur := NewRuleStat(line, f.rulesAggregator)
		f.statsAggregator.addUnbrokenRuleStat(name, *ur)

		f.statsAggregator.sumValuesToAverageUnbroken(ur.Id, ur.Counter, ur.Latency)
	}
}

func (f *FalcoTracer) getBrokenRules() {

	for {
		line := f.falcoGateway.getLine()

		if strings.Contains(string(line), "END BROKEN RULES") {
			break
		}

		name, br := NewRuleStat(line, f.rulesAggregator)
		f.statsAggregator.addBrokenRuleStat(name, *br)

		f.statsAggregator.sumValuesToAverageBroken(br.Id, br.Counter, br.Latency)
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
		Stats StatsAggregator `json:"statistics"`
		Rules RulesAggregator `json:"rules"`
	}{
		Stats: *f.statsAggregator,
		Rules: *f.rulesAggregator,
	})
}
