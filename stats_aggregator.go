package main

type StatsAggregator struct {
	StartTime uint64        `json:"start time"`
	EndTime   uint64        `json:"end time"`
	Fs        []*Falcostats `json:"falco statistics"`
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

func (s *StatsAggregator) setTimes() {
	fsStart := s.Fs[0]
	fsEnd := s.Fs[len(s.Fs)-1]

	s.StartTime = fsStart.StartTime
	s.EndTime = fsEnd.EndTime
}
