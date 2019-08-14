package metrics

type OfflineMetrics struct {
	Fm []*OfflineFalcoMetrics `json:"falco metrics"`
}

func (s *OfflineMetrics) AddOfflineMetrics() {
	f := NewOfflineFalcoMetrics()

	s.Fm = append(s.Fm, f)
}

func (m *OfflineMetrics) AddStackTrace(n string, st Stacktrace) {
	fs := m.Fm[len(m.Fm)-1]

	fs.addStackTrace(n, st)
}

func (m *OfflineMetrics) AddFuncStat(key string, nameFunc string, f FuncStat) {
	fs := m.Fm[len(m.Fm)-1]

	fs.addFuncStat(key, nameFunc, f)
}

func (m *OfflineMetrics) AddCounterStat(n string, cs CounterStat) {
	fs := m.Fm[len(m.Fm)-1]

	fs.addCounterStat(n, cs)
}
