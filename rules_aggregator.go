package main

type RulesAggregator struct {
	falcoRules map[int]FalcoRule
}

func NewRulesAggregator() *RulesAggregator {
	r := new(RulesAggregator)

	r.falcoRules = make(map[int]FalcoRule)

	return r
}

func (r *RulesAggregator) addRule(rule FalcoRule) {
	r.falcoRules[rule.Id] = rule
}

func (r *RulesAggregator) getRuleById(id int) FalcoRule {
	return r.falcoRules[id]
}
