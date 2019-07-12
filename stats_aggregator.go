package main

type StatsAggregator struct {
	Fs []*Falcostats
}

func (s *StatsAggregator) addFalcoStats(start uint64, end uint64) {
	f := new(Falcostats)

	f.StartTime = start
	f.EndTime = end

	s.Fs = append(s.Fs, f)
}

func (s *StatsAggregator) addFuncStat(f FuncStat) {
	fs := s.Fs[len(s.Fs)-1]

	fs.FuncStats = append(fs.FuncStats, f)
}

func (s *StatsAggregator) addCounterStat(cs CounterStat) {
	fs := s.Fs[len(s.Fs)-1]

	fs.CounterStats = append(fs.CounterStats, cs)
}

func (s *StatsAggregator) addUnbrokenRuleStat(r RuleStat) {
	fs := s.Fs[len(s.Fs)-1]

	fs.UnbrokenRules = append(fs.UnbrokenRules, r)
}

func (s *StatsAggregator) addBrokenRuleStat(r RuleStat) {
	fs := s.Fs[len(s.Fs)-1]

	fs.BrokenRules = append(fs.BrokenRules, r)
}
