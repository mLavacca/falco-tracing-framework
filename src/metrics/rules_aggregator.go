package main

type RulesAggregator struct {
	NRules     int               `json:"# of rules"`
	FalcoRules map[int]FalcoRule `json:"falco rules"`
}

func NewRulesAggregator() *RulesAggregator {
	r := new(RulesAggregator)

	r.FalcoRules = make(map[int]FalcoRule)

	return r
}

func (r *RulesAggregator) addRule(rule FalcoRule) {
	r.FalcoRules[rule.Id] = rule
}

func (r *RulesAggregator) getRuleById(id int) FalcoRule {
	return r.FalcoRules[id]
}

func (r *RulesAggregator) setNRules() {
	r.NRules = len(r.FalcoRules)
}
