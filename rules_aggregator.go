package main

type RulesAggregator struct {
	// dummy data structure, to change with a map
	falcoRules []FalcoRule
}

func (r *RulesAggregator) addRule(rule FalcoRule) {
	r.falcoRules = append(r.falcoRules, rule)
}

func (r *RulesAggregator) getRuleNameById(id int) string {

	for i := 0; i < len(r.falcoRules); i++ {
		if r.falcoRules[i].id == id {
			return r.falcoRules[i].name
		}
	}

	return ""
}
