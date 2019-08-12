package metrics

import (
	"sort"
)

type StatsAggregator struct {
	StartTime             uint64        `json:"start time"`
	EndTime               uint64        `json:"end time"`
	Fs                    []*Falcostats `json:"falco statistics"`
	AvgUnbrokenRulesStats []RuleStatAvg `json:"avg unbroken rules stats"`
	AvgBrokenRulesStats   []RuleStatAvg `json:"avg broken rules stats"`
}

func (s *StatsAggregator) AddFalcoStats(start uint64, end uint64) {
	f := NewFalcoStats()

	f.StartTime = start
	f.EndTime = end

	s.Fs = append(s.Fs, f)
}

func (s *StatsAggregator) AddStackTrace(n string, st Stacktrace) {
	fs := s.Fs[len(s.Fs)-1]

	fs.addStackTrace(n, st)
}

func (s *StatsAggregator) AddFuncStat(key string, nameFunc string, f FuncStat) {
	fs := s.Fs[len(s.Fs)-1]

	fs.addFuncStat(key, nameFunc, f)
}

func (s *StatsAggregator) AddCounterStat(n string, cs CounterStat) {
	fs := s.Fs[len(s.Fs)-1]

	fs.addCounterStat(n, cs)
}

func (s *StatsAggregator) AddUnbrokenRuleStat(n string, r RuleStat) {
	fs := s.Fs[len(s.Fs)-1]

	fs.addUnbrokenRule(n, r)
}

func (s *StatsAggregator) AddBrokenRuleStat(n string, r RuleStat) {
	fs := s.Fs[len(s.Fs)-1]

	fs.addBrokenRule(n, r)
}

func (s *StatsAggregator) SumValuesToAverageUnbroken(id int, counter uint64, latency uint64) {
	ol := s.AvgUnbrokenRulesStats[id-1].Rule.Latency
	oc := s.AvgUnbrokenRulesStats[id-1].Rule.Counter

	nl := ((ol * oc) + (latency * counter)) / (oc + counter)

	s.AvgUnbrokenRulesStats[id-1].Rule.Counter += counter
	s.AvgUnbrokenRulesStats[id-1].Rule.Latency = nl
}

func (s *StatsAggregator) SumValuesToAverageBroken(id int, counter uint64, latency uint64) {
	ol := s.AvgBrokenRulesStats[id-1].Rule.Latency
	oc := s.AvgBrokenRulesStats[id-1].Rule.Counter

	nl := ((ol * oc) + (latency * counter)) / (oc + counter)

	s.AvgBrokenRulesStats[id-1].Rule.Counter += counter
	s.AvgBrokenRulesStats[id-1].Rule.Latency = nl
}

func (s *StatsAggregator) SetTimes() {
	fsStart := s.Fs[0]
	fsEnd := s.Fs[len(s.Fs)-1]

	s.StartTime = fsStart.StartTime
	s.EndTime = fsEnd.EndTime
}

func (s *StatsAggregator) SortAvgSlices() {
	sort.SliceStable(s.AvgUnbrokenRulesStats, func(i, j int) bool {
		return s.AvgUnbrokenRulesStats[i].Rule.Latency < s.AvgUnbrokenRulesStats[j].Rule.Latency
	})

	sort.SliceStable(s.AvgBrokenRulesStats, func(i, j int) bool {
		return s.AvgBrokenRulesStats[i].Rule.Latency < s.AvgBrokenRulesStats[j].Rule.Latency
	})
}
