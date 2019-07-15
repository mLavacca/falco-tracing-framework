package main

import "encoding/json"

type StatsAggregator struct {
	Fs []*Falcostats `json:"falco statistics"`
}

func (s *StatsAggregator) addFalcoStats(start uint64, end uint64) {
	f := NewFalcoStats()

	f.StartTime = start
	f.EndTime = end

	s.Fs = append(s.Fs, f)
}

func (s *StatsAggregator) addFuncStat(n string, f FuncStat) {
	fs := s.Fs[len(s.Fs)-1]

	fs.addFuncStat(n, f)
}

func (s *StatsAggregator) addCounterStat(n string, cs CounterStat) {
	fs := s.Fs[len(s.Fs)-1]

	fs.addCounterStat(n, cs)
}

func (s *StatsAggregator) addUnbrokenRuleStat(n string, r RuleStat) {
	fs := s.Fs[len(s.Fs)-1]

	fs.addUnbrokenRule(n, r)
}

func (s *StatsAggregator) addBrokenRuleStat(n string, r RuleStat) {
	fs := s.Fs[len(s.Fs)-1]

	fs.addBrokenRule(n, r)
}

func (s *StatsAggregator) MarshalJSON() ([]byte, error) {
	type StatsAlias StatsAggregator
	return json.Marshal((*StatsAlias)(s))
}
